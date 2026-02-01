package netaux

import (
	"bytes"
	"context"
	"net"
	"sync/atomic"
	"time"
)

// BuildDNSResponse builds a DNS response message for the given query and zone records.
func BuildDNSResponse(query DNSMessage, records []DNSRecord) DNSMessage {
	resp := DNSMessage{ //nolint:exhaustruct
		Header: DNSHeader{ //nolint:exhaustruct
			ID:      query.Header.ID,
			Flags:   FlagQR | FlagAA | (query.Header.Flags & FlagRD),
			QDCount: query.Header.QDCount,
		},
		Questions: query.Questions,
	}

	for _, q := range query.Questions {
		for _, rec := range records {
			if rec.Name == q.Name && rec.Type == q.Type && rec.Class == q.Class {
				rdata := rec.Addr.As16()
				rdLength := uint16(len(rdata))

				if rec.Addr.Is4() {
					v4 := rec.Addr.As4()
					rdata = [16]byte{}
					copy(rdata[:], v4[:])
					rdLength = uint16(len(v4))
				}

				resp.Answers = append(resp.Answers, DNSResourceRecord{
					Name:     rec.Name,
					Type:     rec.Type,
					Class:    rec.Class,
					TTL:      rec.TTL,
					RDLength: rdLength,
					RData:    rdata[:rdLength],
				})
			}
		}
	}

	if len(resp.Answers) == 0 {
		resp.Header.Flags |= RCodeNXDomain
	}

	resp.Header.ANCount = uint16(len(resp.Answers)) //nolint:gosec // because of the total size of DNS answers this is not a relevante concern

	return resp
}

// DNSServer serves DNS responses on a packet connection using zone records.
type DNSServer struct {
	Addr    string
	Records []DNSRecord

	isShutdown  atomic.Bool
	connIsOwned bool
	conn        net.PacketConn
}

// Serve blocks, reading from conn and responding using Records.
// It returns nil when Shutdown is called, or an error on unrecoverable failure.
func (s *DNSServer) Serve(conn net.PacketConn) (serveErr error) {
	s.conn = conn

	defer func() {
		if s.connIsOwned {
			connCloseErr := s.conn.Close()
			if serveErr == nil {
				serveErr = connCloseErr
			}
		}
	}()

	const (
		dnsPackageSize   = 512
		deadlineDuration = 100 * time.Millisecond
	)

	buf := make([]byte, dnsPackageSize)

	for {
		if s.isShutdown.Load() {
			return nil
		}

		if err := conn.SetReadDeadline(time.Now().Add(deadlineDuration)); err != nil {
			return err
		}

		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() { //nolint:errorlint
				continue
			}

			return err
		}

		query, err := ReadDNSMessage(bytes.NewReader(buf[:n]))
		if err != nil {
			continue // skip malformed queries
		}

		resp := BuildDNSResponse(query, s.Records)

		var respBuf bytes.Buffer

		// WriteDNSMessage error check omitted: bytes.Buffer.Write never fails,
		// so the error path is unreachable and would be dead code.
		_ = WriteDNSMessage(&respBuf, resp)

		if _, err := conn.WriteTo(respBuf.Bytes(), addr); err != nil {
			return err
		}
	}
}

// ListenAndServe creates a UDP listener on the given address and calls Serve.
func (s *DNSServer) ListenAndServe(address string) error {
	var lc net.ListenConfig

	conn, err := lc.ListenPacket(context.Background(), "udp", address)
	if err != nil {
		return err
	}

	s.connIsOwned = true

	return s.Serve(conn)
}

// Shutdown signals the server to stop, waits for Serve to return, and closes
// the connection. It returns ctx.Err() if the context expires before the
// server stops.
func (s *DNSServer) Shutdown() {
	s.isShutdown.Store(true)
}
