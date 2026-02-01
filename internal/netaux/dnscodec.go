package netaux

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

// ReadDNSMessage parses a DNS message from the wire format.
func ReadDNSMessage(r io.Reader) (DNSMessage, error) {
	var msg DNSMessage

	if err := binary.Read(r, binary.BigEndian, &msg.Header); err != nil {
		return DNSMessage{}, fmt.Errorf("reading header: %w", err)
	}

	msg.Questions = make([]DNSQuestion, msg.Header.QDCount)
	for i := range msg.Questions {
		q, err := readDNSQuestion(r)
		if err != nil {
			return DNSMessage{}, fmt.Errorf("reading question %d: %w", i, err)
		}

		msg.Questions[i] = q
	}

	msg.Answers = make([]DNSResourceRecord, msg.Header.ANCount)
	for i := range msg.Answers {
		rr, err := readDNSResourceRecord(r)
		if err != nil {
			return DNSMessage{}, fmt.Errorf("reading answer %d: %w", i, err)
		}

		msg.Answers[i] = rr
	}

	msg.Authorities = make([]DNSResourceRecord, msg.Header.NSCount)
	for i := range msg.Authorities {
		rr, err := readDNSResourceRecord(r)
		if err != nil {
			return DNSMessage{}, fmt.Errorf("reading authority %d: %w", i, err)
		}

		msg.Authorities[i] = rr
	}

	msg.Additionals = make([]DNSResourceRecord, msg.Header.ARCount)
	for i := range msg.Additionals {
		rr, err := readDNSResourceRecord(r)
		if err != nil {
			return DNSMessage{}, fmt.Errorf("reading additional %d: %w", i, err)
		}

		msg.Additionals[i] = rr
	}

	return msg, nil
}

// WriteDNSMessage writes a DNS message in wire format.
func WriteDNSMessage(w io.Writer, msg DNSMessage) error {
	if err := binary.Write(w, binary.BigEndian, &msg.Header); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	for i, q := range msg.Questions {
		if err := writeDNSQuestion(w, q); err != nil {
			return fmt.Errorf("writing question %d: %w", i, err)
		}
	}

	for i, rr := range msg.Answers {
		if err := writeDNSResourceRecord(w, rr); err != nil {
			return fmt.Errorf("writing answer %d: %w", i, err)
		}
	}

	for i, rr := range msg.Authorities {
		if err := writeDNSResourceRecord(w, rr); err != nil {
			return fmt.Errorf("writing authority %d: %w", i, err)
		}
	}

	for i, rr := range msg.Additionals {
		if err := writeDNSResourceRecord(w, rr); err != nil {
			return fmt.Errorf("writing additional %d: %w", i, err)
		}
	}

	return nil
}

func readDNSName(r io.Reader) (string, error) {
	var labels []string

	for {
		var length [1]byte
		if _, err := io.ReadFull(r, length[:]); err != nil {
			return "", fmt.Errorf("reading label length: %w", err)
		}

		if length[0] == 0 {
			break
		}

		label := make([]byte, length[0])
		if _, err := io.ReadFull(r, label); err != nil {
			return "", fmt.Errorf("reading label: %w", err)
		}

		labels = append(labels, string(label))
	}

	return strings.Join(labels, ".") + ".", nil
}

func readDNSQuestion(r io.Reader) (DNSQuestion, error) {
	name, err := readDNSName(r)
	if err != nil {
		return DNSQuestion{}, err
	}

	var fields struct {
		Type  uint16
		Class uint16
	}
	if err := binary.Read(r, binary.BigEndian, &fields); err != nil {
		return DNSQuestion{}, fmt.Errorf("reading question fields: %w", err)
	}

	return DNSQuestion{Name: name, Type: fields.Type, Class: fields.Class}, nil
}

func readDNSResourceRecord(r io.Reader) (DNSResourceRecord, error) {
	name, err := readDNSName(r)
	if err != nil {
		return DNSResourceRecord{}, err
	}

	var fields struct {
		Type     uint16
		Class    uint16
		TTL      uint32
		RDLength uint16
	}
	if err := binary.Read(r, binary.BigEndian, &fields); err != nil {
		return DNSResourceRecord{}, fmt.Errorf("reading resource record fields: %w", err)
	}

	rdata := make([]byte, fields.RDLength)
	if _, err := io.ReadFull(r, rdata); err != nil {
		return DNSResourceRecord{}, fmt.Errorf("reading rdata: %w", err)
	}

	return DNSResourceRecord{
		Name:     name,
		Type:     fields.Type,
		Class:    fields.Class,
		TTL:      fields.TTL,
		RDLength: fields.RDLength,
		RData:    rdata,
	}, nil
}

func writeDNSName(w io.Writer, name string) error {
	// Remove trailing dot for splitting
	name = strings.TrimSuffix(name, ".")
	labels := strings.Split(name, ".")

	for _, label := range labels {
		if err := binary.Write(w, binary.BigEndian, uint8(len(label))); err != nil { //nolint:gosec // label length does not exceed 255 bytes.
			return fmt.Errorf("writing label length: %w", err)
		}

		if _, err := w.Write([]byte(label)); err != nil {
			return fmt.Errorf("writing label: %w", err)
		}
	}

	// Write a terminating zero-length label
	if err := binary.Write(w, binary.BigEndian, uint8(0)); err != nil {
		return fmt.Errorf("writing terminator: %w", err)
	}

	return nil
}

func writeDNSQuestion(w io.Writer, q DNSQuestion) error {
	if err := writeDNSName(w, q.Name); err != nil {
		return err
	}

	fields := struct {
		Type  uint16
		Class uint16
	}{q.Type, q.Class}

	return binary.Write(w, binary.BigEndian, &fields)
}

func writeDNSResourceRecord(w io.Writer, rr DNSResourceRecord) error {
	if err := writeDNSName(w, rr.Name); err != nil {
		return err
	}

	fields := struct {
		Type     uint16
		Class    uint16
		TTL      uint32
		RDLength uint16
	}{rr.Type, rr.Class, rr.TTL, rr.RDLength}

	if err := binary.Write(w, binary.BigEndian, &fields); err != nil {
		return err
	}

	if _, err := w.Write(rr.RData); err != nil {
		return fmt.Errorf("writing rdata: %w", err)
	}

	return nil
}
