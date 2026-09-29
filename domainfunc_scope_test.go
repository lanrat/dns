package dns_test

import (
	"strings"
	"testing"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/dnsutil"
)

// TestSetDomainFuncScope checks that the zone parser hook does not leak into the public dnsutil
// package.
func TestSetDomainFuncScope(t *testing.T) {
	defer dns.SetDomainFunc(nil)
	dns.SetDomainFunc(strings.ToLower)
	if got := dnsutil.Absolute("A", "example."); got != "A.example." {
		t.Errorf("dnsutil.Absolute used the zone parser hook: %q", got)
	}
}
