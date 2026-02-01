// Package netaux provides auxiliary utilities for working with network operations.
package netaux

import "net/netip"

// DNS record type constants.
const (
	TypeA    uint16 = 1
	TypeAAAA uint16 = 28
)

// DNS class constants.
const (
	ClassIN uint16 = 1
)

// DNS response code constants.
const (
	RCodeNoError  uint16 = 0
	RCodeNXDomain uint16 = 3
)

// DNS header flag bit positions.
const (
	FlagQR uint16 = 1 << 15 // Query/Response
	FlagAA uint16 = 1 << 10 // Authoritative Answer
	FlagRD uint16 = 1 << 8  // Recursion Desired
	FlagRA uint16 = 1 << 7  // Recursion Available
)

// DNSHeader represents a DNS message header.
type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

// DNSQuestion represents a question entry in a DNS message.
type DNSQuestion struct {
	Name  string // FQDN with trailing dot, e.g. "example.test."
	Type  uint16
	Class uint16
}

// DNSResourceRecord represents a resource record in a DNS message.
type DNSResourceRecord struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte
}

// DNSMessage represents a complete DNS message.
type DNSMessage struct {
	Header      DNSHeader
	Questions   []DNSQuestion
	Answers     []DNSResourceRecord
	Authorities []DNSResourceRecord
	Additionals []DNSResourceRecord
}

// DNSRecord represents a user-facing DNS record for zone data.
type DNSRecord struct {
	Name  string     // FQDN with trailing dot
	Type  uint16
	Class uint16
	TTL   uint32
	Addr  netip.Addr // works for IPv4 and IPv6
}
