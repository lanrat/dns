package dns

import (
	"strings"
	"testing"
)

// TestZoneParserCaseInsensitive checks that classes, types, TYPExxx and CLASSxxx are accepted in any
// case, as RFC 1035 Section 5.1 requires. Some published zone files are entirely lower case.
func TestZoneParserCaseInsensitive(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"a.example. 3600 in a 127.0.0.1", "a.example.\t3600\tIN\tA\t127.0.0.1"},
		{"a.example. in 3600 ns ns.example.", "a.example.\t3600\tIN\tNS\tns.example."},
		{"a.example. 3600 class1 a 127.0.0.1", "a.example.\t3600\tIN\tA\t127.0.0.1"},
		{"a.example. 3600 in type65534 \\# 1 00", "a.example.\t3600\tIN\tTYPE65534\t\\# 1 00"},
		{
			"a.example. 3600 in rrsig soa 8 2 3600 20260101000000 20250101000000 1 example. AAAA",
			"a.example.\t3600\tIN\tRRSIG\tSOA 8 2 3600 20260101000000 20250101000000 1 example. AAAA",
		},
		{
			"a.example. 3600 in rrsig type65534 8 2 3600 20260101000000 20250101000000 1 example. AAAA",
			"a.example.\t3600\tIN\tRRSIG\tTYPE65534 8 2 3600 20260101000000 20250101000000 1 example. AAAA",
		},
		{"a.example. 3600 in nsec b.example. a ns rrsig type65534", "a.example.\t3600\tIN\tNSEC\tb.example. A NS RRSIG TYPE65534"},
		{"_dsync.example. 3600 in dsync cds notify 5359 t.example.", "_dsync.example.\t3600\tIN\tDSYNC\tCDS NOTIFY 5359 t.example."},
	}
	for _, tc := range tests {
		zp := NewZoneParser(strings.NewReader(tc.in), "", "")
		rr, ok := zp.Next()
		if !ok {
			t.Errorf("%q: %v", tc.in, zp.Err())
			continue
		}
		if got := rr.String(); got != tc.want {
			t.Errorf("%q:\n got  %q\n want %q", tc.in, got, tc.want)
		}
	}
}
