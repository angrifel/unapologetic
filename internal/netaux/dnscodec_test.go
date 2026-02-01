package netaux

import (
	"bytes"
	"errors"
	"slices"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
	"github.com/angrifel/unapologetic/internal/require"
	"github.com/angrifel/unapologetic/iospy"
)

func TestWriteAndReadDNSMessage_QueryWithAQuestion(t *testing.T) {
	original := DNSMessage{
		Header: DNSHeader{
			ID:      0x1234,
			Flags:   FlagRD,
			QDCount: 1,
		},
		Questions: []DNSQuestion{
			{Name: "example.test.", Type: TypeA, Class: ClassIN},
		},
	}

	var buf bytes.Buffer
	err := WriteDNSMessage(&buf, original)
	assert.IsNil(t, err, "WriteDNSMessage should not error")

	parsed, err := ReadDNSMessage(&buf)
	assert.IsNil(t, err, "ReadDNSMessage should not error")

	assert.Equal(t, original.Header, parsed.Header, "headers should match")
	assert.Equal(t, len(original.Questions), len(parsed.Questions), "question count should match")
	assert.Equal(t, original.Questions[0], parsed.Questions[0], "question should match")
}

func TestWriteAndReadDNSMessage_QueryWithAAAAQuestion(t *testing.T) {
	original := DNSMessage{
		Header: DNSHeader{
			ID:      0xABCD,
			Flags:   FlagRD,
			QDCount: 1,
		},
		Questions: []DNSQuestion{
			{Name: "ipv6.test.", Type: TypeAAAA, Class: ClassIN},
		},
	}

	var buf bytes.Buffer
	err := WriteDNSMessage(&buf, original)
	assert.IsNil(t, err, "WriteDNSMessage should not error")

	parsed, err := ReadDNSMessage(&buf)
	assert.IsNil(t, err, "ReadDNSMessage should not error")

	assert.Equal(t, original.Header, parsed.Header, "headers should match")
	assert.Equal(t, len(original.Questions), len(parsed.Questions), "question count should match")
	assert.Equal(t, original.Questions[0], parsed.Questions[0], "question should match")
}

func TestWriteAndReadDNSMessage_ResponseWithAnswers(t *testing.T) {
	original := DNSMessage{
		Header: DNSHeader{
			ID:      0x5678,
			Flags:   FlagQR | FlagAA | FlagRD,
			QDCount: 1,
			ANCount: 1,
		},
		Questions: []DNSQuestion{
			{Name: "example.test.", Type: TypeA, Class: ClassIN},
		},
		Answers: []DNSResourceRecord{
			{
				Name:     "example.test.",
				Type:     TypeA,
				Class:    ClassIN,
				TTL:      300,
				RDLength: 4,
				RData:    []byte{192, 168, 1, 1},
			},
		},
	}

	var buf bytes.Buffer
	err := WriteDNSMessage(&buf, original)
	assert.IsNil(t, err, "WriteDNSMessage should not error")

	parsed, err := ReadDNSMessage(&buf)
	assert.IsNil(t, err, "ReadDNSMessage should not error")

	assert.Equal(t, original.Header, parsed.Header, "headers should match")
	assert.Equal(t, len(original.Questions), len(parsed.Questions), "question count should match")
	assert.Equal(t, original.Questions[0], parsed.Questions[0], "question should match")
	assert.Equal(t, len(original.Answers), len(parsed.Answers), "answer count should match")
	assert.Equal(t, original.Answers[0].Name, parsed.Answers[0].Name, "answer name should match")
	assert.Equal(t, original.Answers[0].Type, parsed.Answers[0].Type, "answer type should match")
	assert.Equal(t, original.Answers[0].Class, parsed.Answers[0].Class, "answer class should match")
	assert.Equal(t, original.Answers[0].TTL, parsed.Answers[0].TTL, "answer TTL should match")
	assert.Equal(t, original.Answers[0].RDLength, parsed.Answers[0].RDLength, "answer RDLength should match")
	assert.Equal(t, true, slices.Equal(original.Answers[0].RData, parsed.Answers[0].RData), "answer RData should match")
}

func TestWriteAndReadDNSMessage_EmptyQuery(t *testing.T) {
	original := DNSMessage{
		Header: DNSHeader{
			ID:    0x0001,
			Flags: 0,
		},
	}

	var buf bytes.Buffer
	err := WriteDNSMessage(&buf, original)
	assert.IsNil(t, err, "WriteDNSMessage should not error")

	parsed, err := ReadDNSMessage(&buf)
	assert.IsNil(t, err, "ReadDNSMessage should not error")

	assert.Equal(t, original.Header, parsed.Header, "headers should match")
	assert.Equal(t, 0, len(parsed.Questions), "should have no questions")
	assert.Equal(t, 0, len(parsed.Answers), "should have no answers")
	assert.Equal(t, 0, len(parsed.Authorities), "should have no authorities")
	assert.Equal(t, 0, len(parsed.Additionals), "should have no additionals")
}

func TestWriteAndReadDNSMessage_ResponseWithMultipleAnswers(t *testing.T) {
	original := DNSMessage{
		Header: DNSHeader{
			ID:      0x9999,
			Flags:   FlagQR | FlagAA,
			QDCount: 1,
			ANCount: 2,
		},
		Questions: []DNSQuestion{
			{Name: "multi.test.", Type: TypeA, Class: ClassIN},
		},
		Answers: []DNSResourceRecord{
			{
				Name:     "multi.test.",
				Type:     TypeA,
				Class:    ClassIN,
				TTL:      60,
				RDLength: 4,
				RData:    []byte{10, 0, 0, 1},
			},
			{
				Name:     "multi.test.",
				Type:     TypeA,
				Class:    ClassIN,
				TTL:      60,
				RDLength: 4,
				RData:    []byte{10, 0, 0, 2},
			},
		},
	}

	var buf bytes.Buffer
	err := WriteDNSMessage(&buf, original)
	assert.IsNil(t, err, "WriteDNSMessage should not error")

	parsed, err := ReadDNSMessage(&buf)
	assert.IsNil(t, err, "ReadDNSMessage should not error")

	assert.Equal(t, original.Header, parsed.Header, "headers should match")
	assert.Equal(t, 2, len(parsed.Answers), "should have two answers")
	assert.Equal(t, true, slices.Equal(original.Answers[0].RData, parsed.Answers[0].RData), "first answer RData should match")
	assert.Equal(t, true, slices.Equal(original.Answers[1].RData, parsed.Answers[1].RData), "second answer RData should match")
}

func TestWriteAndReadDNSMessage_AllSections(t *testing.T) {
	original := DNSMessage{
		Header: DNSHeader{
			ID:      0xFACE,
			Flags:   FlagQR | FlagAA,
			QDCount: 1,
			ANCount: 1,
			NSCount: 1,
			ARCount: 1,
		},
		Questions: []DNSQuestion{
			{Name: "ns.test.", Type: TypeA, Class: ClassIN},
		},
		Answers: []DNSResourceRecord{
			{Name: "ns.test.", Type: TypeA, Class: ClassIN, TTL: 300, RDLength: 4, RData: []byte{10, 0, 0, 1}},
		},
		Authorities: []DNSResourceRecord{
			{Name: "ns.test.", Type: TypeA, Class: ClassIN, TTL: 300, RDLength: 4, RData: []byte{10, 0, 0, 2}},
		},
		Additionals: []DNSResourceRecord{
			{Name: "ns.test.", Type: TypeA, Class: ClassIN, TTL: 300, RDLength: 4, RData: []byte{10, 0, 0, 3}},
		},
	}

	var buf bytes.Buffer
	err := WriteDNSMessage(&buf, original)
	assert.IsNil(t, err, "WriteDNSMessage should not error")

	parsed, err := ReadDNSMessage(&buf)
	assert.IsNil(t, err, "ReadDNSMessage should not error")

	assert.Equal(t, original.Header, parsed.Header, "headers should match")
	assert.Equal(t, 1, len(parsed.Questions), "should have one question")
	assert.Equal(t, original.Questions[0], parsed.Questions[0], "question should match")
	assert.Equal(t, 1, len(parsed.Answers), "should have one answer")
	assert.Equal(t, true, slices.Equal(original.Answers[0].RData, parsed.Answers[0].RData), "answer RData should match")
	assert.Equal(t, 1, len(parsed.Authorities), "should have one authority")
	assert.Equal(t, true, slices.Equal(original.Authorities[0].RData, parsed.Authorities[0].RData), "authority RData should match")
	assert.Equal(t, 1, len(parsed.Additionals), "should have one additional")
	assert.Equal(t, true, slices.Equal(original.Additionals[0].RData, parsed.Additionals[0].RData), "additional RData should match")
}

// fullDNSMessage returns a DNS message with all sections populated (question, answer,
// authority, additional) using "ns.test." as the name. This produces a 94-byte wire
// representation with the following byte layout:
//
//	Bytes  0-11: Header (12 bytes)
//	Bytes 12-24: Question (name 12-20, fields 21-24)
//	Bytes 25-47: Answer RR (name 25-33, fields 34-43, rdata 44-47)
//	Bytes 48-70: Authority RR (name 48-56, fields 57-66, rdata 67-70)
//	Bytes 71-93: Additional RR (name 71-79, fields 80-89, rdata 90-93)
func fullDNSMessage() DNSMessage {
	return DNSMessage{
		Header: DNSHeader{
			ID:      0xFACE,
			Flags:   FlagQR | FlagAA,
			QDCount: 1,
			ANCount: 1,
			NSCount: 1,
			ARCount: 1,
		},
		Questions: []DNSQuestion{
			{Name: "ns.test.", Type: TypeA, Class: ClassIN},
		},
		Answers: []DNSResourceRecord{
			{Name: "ns.test.", Type: TypeA, Class: ClassIN, TTL: 300, RDLength: 4, RData: []byte{10, 0, 0, 1}},
		},
		Authorities: []DNSResourceRecord{
			{Name: "ns.test.", Type: TypeA, Class: ClassIN, TTL: 300, RDLength: 4, RData: []byte{10, 0, 0, 2}},
		},
		Additionals: []DNSResourceRecord{
			{Name: "ns.test.", Type: TypeA, Class: ClassIN, TTL: 300, RDLength: 4, RData: []byte{10, 0, 0, 3}},
		},
	}
}

func TestReadDNSMessage_TruncatedInput(t *testing.T) {
	// Serialize the full message to get valid wire bytes, then truncate at various offsets
	// to exercise each error branch in ReadDNSMessage and its sub-functions.
	msg := fullDNSMessage()
	var fullBuf bytes.Buffer
	err := WriteDNSMessage(&fullBuf, msg)
	assert.IsNil(t, err, "WriteDNSMessage should not error")
	wireBytes := fullBuf.Bytes()

	testCases := []struct {
		name       string
		truncateAt int
	}{
		{"truncated header", 6},
		{"truncated question name at label length", 12},
		{"truncated question name at label data", 14},
		{"truncated question fields", 22},
		{"truncated answer name", 26},
		{"truncated answer fields", 36},
		{"truncated answer rdata", 45},
		{"truncated authority name", 49},
		{"truncated authority fields", 59},
		{"truncated authority rdata", 68},
		{"truncated additional name", 72},
		{"truncated additional fields", 82},
		{"truncated additional rdata", 91},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ReadDNSMessage(bytes.NewReader(wireBytes[:tc.truncateAt]))
			assert.IsNotNil(t, err, "ReadDNSMessage should error on truncated input")
		})
	}
}

func TestWriteDNSMessage_WriterError(t *testing.T) {
	// Serialize the full message to determine the total wire size, then use
	// LimitWriterWithError to trigger write failures at specific byte offsets.
	// This covers every error branch in WriteDNSMessage and its sub-functions.
	//
	// Using iospy.LimitWriterWithError because bytes.Buffer.Write never fails,
	// making it impossible to exercise write error paths with real I/O alone.
	msg := fullDNSMessage()
	var fullBuf bytes.Buffer
	err := WriteDNSMessage(&fullBuf, msg)
	assert.IsNil(t, err, "WriteDNSMessage should not error")
	totalSize := int64(fullBuf.Len())

	writeErr := errors.New("write limit reached")

	testCases := []struct {
		name    string
		limitAt int64
	}{
		{"error writing header", 6},
		{"error writing question name label length", 12},
		{"error writing question name label data", 14},
		{"error writing question name terminator", 20},
		{"error writing question fields", 22},
		{"error writing answer name label length", 25},
		{"error writing answer name label data", 27},
		{"error writing answer name terminator", 33},
		{"error writing answer fields", 36},
		{"error writing answer rdata", 45},
		{"error writing authority name", 48},
		{"error writing authority fields", 59},
		{"error writing authority rdata", 68},
		{"error writing additional name", 71},
		{"error writing additional fields", 82},
		{"error writing additional rdata", 91},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Less(t, tc.limitAt, totalSize)
			w := iospy.LimitWriterWithError(&bytes.Buffer{}, tc.limitAt, writeErr)
			err := WriteDNSMessage(w, msg)
			assert.IsNotNil(t, err, "WriteDNSMessage should error when writer fails")
		})
	}
}
