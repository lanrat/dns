package dnslex

import "strings"

// LookupFold looks up s in m ignoring case. Mnemonics in zone files are case-insensitive (RFC 1035,
// Section 5.1), and the maps are keyed by the upper case form.
func LookupFold(s string, m map[string]uint16) (uint16, bool) { return upperLookup(s, m) }

// HasPrefixFold reports whether s begins with prefix, ignoring case.
func HasPrefixFold(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}
