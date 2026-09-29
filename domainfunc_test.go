package dns

import (
	"strings"
	"testing"
)

// longLabel is 76 bytes in UTF-8, too long for a DNS label, while toPunycode below turns it into a
// valid 42 byte label.
const longLabel = "ที่ปรึกษาธุรกิจ-แผนการตลาด"

// longLabelACE is the punycode form of longLabel.
const longLabelACE = "xn----twfab7a0egkknqr9fcg7a9c1jid6anz6c8o7f"

// toPunycode is a stand-in for a real IDNA conversion: it lower-cases and replaces longLabel.
func toPunycode(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), longLabel, longLabelACE)
}

func TestSetDomainFunc(t *testing.T) {
	tests := []struct {
		zone   string
		origin string
		want   string // owner and rdata of the parsed RR, or "" when parsing must fail
		hook   bool
	}{
		{"A.Example. 3600 IN NS NS1.Example.", "", "a.example. ns1.example.", true},
		{longLabel + ".example. 3600 IN NS ns1.example.", "", longLabelACE + ".example. ns1.example.", true},
		{longLabel + " 3600 IN NS ns1.example.", "example.", longLabelACE + ".example. ns1.example.", true},
		{"a.example. 3600 IN NS ns." + longLabel + ".example.", "", "a.example. ns." + longLabelACE + ".example.", true},
		{"a 3600 IN NS ns1.example.", longLabel + ".example.", "a." + longLabelACE + ".example. ns1.example.", true},
		{longLabel + ".example. 3600 IN NS ns1.example.", "", "", false},
		{"a.example. 3600 IN NS ns." + longLabel + ".example.", "", "", false},
	}
	defer SetDomainFunc(nil)
	for _, tc := range tests {
		if tc.hook {
			SetDomainFunc(toPunycode)
		} else {
			SetDomainFunc(nil)
		}
		zp := NewZoneParser(strings.NewReader(tc.zone), tc.origin, "")
		rr, ok := zp.Next()
		if tc.want == "" {
			if ok {
				t.Errorf("%q without hook: expected an error, got %s", tc.zone, rr)
			}
			continue
		}
		if !ok {
			t.Errorf("%q: %v", tc.zone, zp.Err())
			continue
		}
		if got := rr.Header().Name + " " + rr.(*NS).Ns; got != tc.want {
			t.Errorf("%q: got %q, want %q", tc.zone, got, tc.want)
		}
	}
}
