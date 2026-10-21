// Package dnsnet has the wire format reader used while unpacking.
package dnsnet

import "golang.org/x/crypto/cryptobyte"

// Pointer follows a message may spend, per octet and at least.
const (
	followsPerOctet = 2
	minFollows      = 2048
)

// String is a cryptobyte.String that carries the compression pointer follows its message has left.
type String struct {
	cryptobyte.String
	left uint32
}

// New returns a String over msg with a follow budget proportional to its length.
func New(msg []byte) String {
	return String{String: msg, left: uint32(max(followsPerOctet*len(msg), minFollows))}
}

// Follow spends one pointer follow and reports whether the budget allowed it.
func (s *String) Follow() bool {
	if s.left == 0 {
		return false
	}
	s.left--
	return true
}
