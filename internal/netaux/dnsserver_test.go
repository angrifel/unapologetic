package netaux

import (
	"bytes"
	"context"
	"errors"
	"net"
	"net/netip"
	"sync"
	"testing"
	"time"

	"github.com/angrifel/unapologetic/internal/assert"
)

func TestBuildDNSResponse_MatchingARecord(t *testing.T) {
	query := DNSMessage{
		Header: DNSHeader{
			ID:      0x1234,
			Flags:   FlagRD,
			QDCount: 1,
		},
		Questions: []DNSQuestion{
			{Name: "example.test.", Type: TypeA, Class: ClassIN},
		},
	}

	records := []DNSRecord{
		{Name: "example.test.", Type: TypeA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("192.168.1.1")},
	}

	resp := BuildDNSResponse(query, records)

	assert.Equal(t, query.Header.ID, resp.Header.ID, "ID should be preserved")
	assert.Equal(t, true, resp.Header.Flags&FlagQR != 0, "QR flag should be set")
	assert.Equal(t, true, resp.Header.Flags&FlagAA != 0, "AA flag should be set")
	assert.Equal(t, true, resp.Header.Flags&FlagRD != 0, "RD flag should be copied from query")
	assert.Equal(t, uint16(1), resp.Header.ANCount, "should have one answer")
	assert.Equal(t, "example.test.", resp.Answers[0].Name, "answer name should match")
	assert.Equal(t, TypeA, resp.Answers[0].Type, "answer type should be A")
	assert.Equal(t, uint16(4), resp.Answers[0].RDLength, "A record should be 4 bytes")
	assert.Equal(t, byte(192), resp.Answers[0].RData[0], "first octet")
	assert.Equal(t, byte(168), resp.Answers[0].RData[1], "second octet")
	assert.Equal(t, byte(1), resp.Answers[0].RData[2], "third octet")
	assert.Equal(t, byte(1), resp.Answers[0].RData[3], "fourth octet")
}

func TestBuildDNSResponse_NXDomain(t *testing.T) {
	query := DNSMessage{
		Header: DNSHeader{
			ID:      0x5678,
			Flags:   FlagRD,
			QDCount: 1,
		},
		Questions: []DNSQuestion{
			{Name: "nonexistent.test.", Type: TypeA, Class: ClassIN},
		},
	}

	resp := BuildDNSResponse(query, nil)

	assert.Equal(t, query.Header.ID, resp.Header.ID, "ID should be preserved")
	assert.Equal(t, RCodeNXDomain, resp.Header.Flags&0x000F, "should return NXDomain")
	assert.Equal(t, uint16(0), resp.Header.ANCount, "should have no answers")
}

func TestBuildDNSResponse_IDPreservation(t *testing.T) {
	query := DNSMessage{
		Header: DNSHeader{
			ID:      0xBEEF,
			QDCount: 1,
		},
		Questions: []DNSQuestion{
			{Name: "example.test.", Type: TypeA, Class: ClassIN},
		},
	}

	resp := BuildDNSResponse(query, nil)
	assert.Equal(t, uint16(0xBEEF), resp.Header.ID, "ID should be preserved")
}

func TestBuildDNSResponse_TypeFiltering(t *testing.T) {
	query := DNSMessage{
		Header: DNSHeader{
			ID:      0x0001,
			QDCount: 1,
		},
		Questions: []DNSQuestion{
			{Name: "example.test.", Type: TypeAAAA, Class: ClassIN},
		},
	}

	records := []DNSRecord{
		{Name: "example.test.", Type: TypeA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("192.168.1.1")},
		{Name: "example.test.", Type: TypeAAAA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("::1")},
	}

	resp := BuildDNSResponse(query, records)

	assert.Equal(t, uint16(1), resp.Header.ANCount, "should have one answer (AAAA only)")
	assert.Equal(t, TypeAAAA, resp.Answers[0].Type, "answer should be AAAA")
	assert.Equal(t, uint16(16), resp.Answers[0].RDLength, "AAAA record should be 16 bytes")
}

func TestBuildDNSResponse_MultipleMatchingRecords(t *testing.T) {
	query := DNSMessage{
		Header: DNSHeader{
			ID:      0x0002,
			QDCount: 1,
		},
		Questions: []DNSQuestion{
			{Name: "multi.test.", Type: TypeA, Class: ClassIN},
		},
	}

	records := []DNSRecord{
		{Name: "multi.test.", Type: TypeA, Class: ClassIN, TTL: 60, Addr: netip.MustParseAddr("10.0.0.1")},
		{Name: "multi.test.", Type: TypeA, Class: ClassIN, TTL: 60, Addr: netip.MustParseAddr("10.0.0.2")},
	}

	resp := BuildDNSResponse(query, records)

	assert.Equal(t, uint16(2), resp.Header.ANCount, "should have two answers")
}

func startDNSServer(t *testing.T, records []DNSRecord) (*DNSServer, net.Addr) {
	t.Helper()
	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	assert.IsNil(t, err, "ListenPacket should not error")
	addr := conn.LocalAddr()
	srv := &DNSServer{Records: records}
	done := make(chan struct{})
	go func() {
		defer close(done)
		srv.Serve(conn)
	}()
	t.Cleanup(func() {
		srv.Shutdown()
		conn.Close()
		<-done
	})
	return srv, addr
}

func TestDNSServer_Serve_LookupHost(t *testing.T) {
	records := []DNSRecord{
		{Name: "myhost.test.", Type: TypeA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("127.0.0.99")},
	}

	_, addr := startDNSServer(t, records)
	resolver := newTestResolver(addr)

	addrs, err := resolver.LookupHost(context.Background(), "myhost.test")
	assert.IsNil(t, err, "LookupHost should not error")
	assert.Equal(t, 1, len(addrs), "should resolve one address")
	assert.Equal(t, "127.0.0.99", addrs[0], "should resolve to expected address")
}

func TestDNSServer_Serve_LookupHostMultipleARecords(t *testing.T) {
	records := []DNSRecord{
		{Name: "multi.test.", Type: TypeA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("10.0.0.1")},
		{Name: "multi.test.", Type: TypeA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("10.0.0.2")},
	}

	_, addr := startDNSServer(t, records)
	resolver := newTestResolver(addr)

	addrs, err := resolver.LookupHost(context.Background(), "multi.test")
	assert.IsNil(t, err, "LookupHost should not error")
	assert.Equal(t, 2, len(addrs), "should resolve two addresses")
}

func TestDNSServer_Serve_LookupIPv4(t *testing.T) {
	records := []DNSRecord{
		{Name: "v4.test.", Type: TypeA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("192.168.1.100")},
	}

	_, addr := startDNSServer(t, records)
	resolver := newTestResolver(addr)

	ips, err := resolver.LookupIP(context.Background(), "ip4", "v4.test")
	assert.IsNil(t, err, "LookupIP should not error")
	assert.Equal(t, 1, len(ips), "should resolve one IP")
	assert.Equal(t, "192.168.1.100", ips[0].String(), "should resolve to expected IPv4")
}

func TestDNSServer_Serve_LookupIPv6(t *testing.T) {
	records := []DNSRecord{
		{Name: "v6.test.", Type: TypeAAAA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("::1")},
	}

	_, addr := startDNSServer(t, records)
	resolver := newTestResolver(addr)

	ips, err := resolver.LookupIP(context.Background(), "ip6", "v6.test")
	assert.IsNil(t, err, "LookupIP should not error")
	assert.Equal(t, 1, len(ips), "should resolve one IP")
	assert.Equal(t, "::1", ips[0].String(), "should resolve to expected IPv6")
}

func TestDNSServer_Serve_NonExistentDomain(t *testing.T) {
	records := []DNSRecord{
		{Name: "exists.test.", Type: TypeA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("10.0.0.1")},
	}

	_, addr := startDNSServer(t, records)
	resolver := newTestResolver(addr)

	_, err := resolver.LookupHost(context.Background(), "doesnotexist.test")
	assert.IsNotNil(t, err, "LookupHost should return error for non-existent domain")
}

func TestDNSServer_Serve_ConcurrentLookups(t *testing.T) {
	records := []DNSRecord{
		{Name: "concurrent.test.", Type: TypeA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("10.1.1.1")},
	}

	_, addr := startDNSServer(t, records)
	resolver := newTestResolver(addr)

	var wg sync.WaitGroup
	const concurrency = 10
	errs := make([]error, concurrency)
	results := make([][]string, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx], errs[idx] = resolver.LookupHost(context.Background(), "concurrent.test")
		}(i)
	}

	wg.Wait()

	for i := 0; i < concurrency; i++ {
		assert.IsNil(t, errs[i], "concurrent lookup %d should not error", i)
		assert.Equal(t, 1, len(results[i]), "concurrent lookup %d should resolve one address", i)
		assert.Equal(t, "10.1.1.1", results[i][0], "concurrent lookup %d should resolve correctly", i)
	}
}

func TestDNSServer_Serve_ShutdownReturnsNil(t *testing.T) {
	records := []DNSRecord{
		{Name: "test.test.", Type: TypeA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("10.0.0.1")},
	}

	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	assert.IsNil(t, err, "ListenPacket should not error")
	defer conn.Close()

	srv := &DNSServer{Records: records}
	done := make(chan error, 1)
	go func() {
		done <- srv.Serve(conn)
	}()

	srv.Shutdown()

	serveErr := <-done
	assert.IsNil(t, serveErr, "Serve should return nil on normal shutdown")
}

func TestDNSServer_Serve_SetReadDeadlineError(t *testing.T) {
	deadlineErr := errors.New("set deadline failed")
	conn := &failingPacketConn{
		setReadDeadlineErr: deadlineErr,
	}

	srv := &DNSServer{}
	err := srv.Serve(conn)
	assert.Equal(t, deadlineErr, err, "Serve should return SetReadDeadline error")
}

func TestDNSServer_Serve_ReadFromNonTimeoutError(t *testing.T) {
	readErr := errors.New("connection reset")
	realConn, err := net.ListenPacket("udp", "127.0.0.1:0")
	assert.IsNil(t, err, "ListenPacket should not error")
	defer realConn.Close()

	conn := &failingPacketConn{
		PacketConn:  realConn,
		readFromErr: readErr,
	}

	srv := &DNSServer{}
	err = srv.Serve(conn)
	assert.Equal(t, readErr, err, "Serve should return ReadFrom error")
}

func TestDNSServer_Serve_MalformedQuery(t *testing.T) {
	records := []DNSRecord{
		{Name: "test.test.", Type: TypeA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("10.0.0.1")},
	}

	_, addr := startDNSServer(t, records)

	// Send malformed data (too short to be a valid DNS message)
	clientConn, err := net.DialTimeout("udp", addr.String(), time.Second)
	assert.IsNil(t, err, "Dial should not error")
	defer clientConn.Close()

	_, err = clientConn.Write([]byte{0xDE, 0xAD})
	assert.IsNil(t, err, "Write should not error")

	// Send a valid query afterward to prove the server is still running
	// (i.e., the malformed query was skipped, not fatal)
	resolver := newTestResolver(addr)
	addrs, err := resolver.LookupHost(context.Background(), "test.test")
	assert.IsNil(t, err, "LookupHost should succeed after malformed query")
	assert.Equal(t, 1, len(addrs), "should resolve one address")
	assert.Equal(t, "10.0.0.1", addrs[0], "should resolve correctly")
}

func TestDNSServer_Serve_WriteToError(t *testing.T) {
	writeErr := errors.New("write failed")
	realConn, err := net.ListenPacket("udp", "127.0.0.1:0")
	assert.IsNil(t, err, "ListenPacket should not error")
	defer realConn.Close()

	conn := &failingPacketConn{
		PacketConn: realConn,
		writeToErr: writeErr,
	}

	srv := &DNSServer{
		Records: []DNSRecord{
			{Name: "test.test.", Type: TypeA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("10.0.0.1")},
		},
	}

	done := make(chan error, 1)
	go func() {
		done <- srv.Serve(conn)
	}()

	// Send a valid DNS query to the server
	clientConn, err := net.DialTimeout("udp", realConn.LocalAddr().String(), time.Second)
	assert.IsNil(t, err, "Dial should not error")
	defer clientConn.Close()

	query := DNSMessage{
		Header:    DNSHeader{ID: 0x1234, Flags: FlagRD, QDCount: 1},
		Questions: []DNSQuestion{{Name: "test.test.", Type: TypeA, Class: ClassIN}},
	}
	var buf bytes.Buffer
	err = WriteDNSMessage(&buf, query)
	assert.IsNil(t, err, "WriteDNSMessage should not error")
	_, err = clientConn.Write(buf.Bytes())
	assert.IsNil(t, err, "Write should not error")

	serveErr := <-done
	assert.Equal(t, writeErr, serveErr, "Serve should return WriteTo error")
}

func TestDNSServer_ListenAndServe_Success(t *testing.T) {
	srv := &DNSServer{Records: []DNSRecord{
		{Name: "test.test.", Type: TypeA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("10.0.0.1")},
	}}

	done := make(chan error, 1)
	go func() {
		done <- srv.ListenAndServe("127.0.0.1:0")
	}()

	// Allow time for the server to bind and enter the serve loop
	time.Sleep(50 * time.Millisecond)

	srv.Shutdown()

	serveErr := <-done
	assert.IsNil(t, serveErr, "ListenAndServe should return nil on shutdown")
}

func TestDNSServer_ListenAndServe_InvalidAddress(t *testing.T) {
	srv := &DNSServer{}
	err := srv.ListenAndServe("invalid-address-no-port")
	assert.IsNotNil(t, err, "ListenAndServe should error on invalid address")
}

func TestDNSServer_ListenAndServe_ClosesConnOnShutdown(t *testing.T) {
	srv := &DNSServer{Records: []DNSRecord{
		{Name: "test.test.", Type: TypeA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("10.0.0.1")},
	}}

	done := make(chan error, 1)
	go func() {
		done <- srv.ListenAndServe("127.0.0.1:0")
	}()

	// Allow time for the server to bind and enter the serve loop
	time.Sleep(50 * time.Millisecond)

	srv.Shutdown()
	<-done

	// The conn should have been closed by Serve's defer since connIsOwned is true.
	// Writing to a closed conn returns an error.
	_, err := srv.conn.WriteTo([]byte{0x00}, srv.conn.LocalAddr())
	assert.IsNotNil(t, err, "conn should be closed after ListenAndServe returns from shutdown")
}

func TestDNSServer_Serve_DoesNotCloseCallerConn(t *testing.T) {
	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	assert.IsNil(t, err, "ListenPacket should not error")
	defer conn.Close()

	srv := &DNSServer{Records: []DNSRecord{
		{Name: "test.test.", Type: TypeA, Class: ClassIN, TTL: 300, Addr: netip.MustParseAddr("10.0.0.1")},
	}}

	done := make(chan error, 1)
	go func() {
		done <- srv.Serve(conn)
	}()

	srv.Shutdown()
	<-done

	// The conn should still be usable since connIsOwned is false.
	err = conn.SetReadDeadline(time.Now().Add(time.Millisecond))
	assert.IsNil(t, err, "conn should still be usable after Serve returns when caller owns it")
}

func TestDNSServer_Serve_ClosesOwnedConnOnError(t *testing.T) {
	// When connIsOwned is true, the conn should be closed even when Serve
	// returns an error. The original serve error must be preserved.
	deadlineErr := errors.New("set deadline failed")
	conn := &failingPacketConn{
		setReadDeadlineErr: deadlineErr,
	}

	srv := &DNSServer{}
	srv.connIsOwned = true

	err := srv.Serve(conn)
	assert.Equal(t, deadlineErr, err, "Serve should return the original error, not the close error")
	assert.Equal(t, true, conn.closed, "conn should be closed even when Serve returns an error")
}

type failingPacketConn struct {
	net.PacketConn
	setReadDeadlineErr error
	readFromErr        error
	writeToErr         error
	closed             bool
}

func (c *failingPacketConn) Close() error {
	c.closed = true
	if c.PacketConn != nil {
		return c.PacketConn.Close()
	}
	return nil
}

func (c *failingPacketConn) ReadFrom(p []byte) (int, net.Addr, error) {
	if c.readFromErr != nil {
		return 0, nil, c.readFromErr
	}
	return c.PacketConn.ReadFrom(p)
}

func (c *failingPacketConn) SetReadDeadline(t time.Time) error {
	if c.setReadDeadlineErr != nil {
		return c.setReadDeadlineErr
	}
	if c.PacketConn != nil {
		return c.PacketConn.SetReadDeadline(t)
	}
	return nil
}

func (c *failingPacketConn) WriteTo(p []byte, addr net.Addr) (int, error) {
	if c.writeToErr != nil {
		return 0, c.writeToErr
	}
	if c.PacketConn != nil {
		return c.PacketConn.WriteTo(p, addr)
	}
	return len(p), nil
}

func newTestResolver(addr net.Addr) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "udp", addr.String())
		},
	}
}
