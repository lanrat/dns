package rdata_test

import (
	"slices"
	"testing"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/rdata"
)

func TestTXT(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		txts []string
		want string
	}{
		{"simple", []string{`simple`}, `"simple"`},
		{"space", []string{`with space`}, `"with space"`},
		{"two", []string{`one`, `two`}, `"one" "two"`},
		{"quote", []string{`qu"ote`}, `"qu\"ote"`},
		{"etx", []string{"\x03"}, `"\003"`},
		{"del", []string{"\x7f"}, `"\127"`},
		{"swearing", []string{`!@#$%^&*();:'"<>,./?~`}, `"!@#$%^&*();:'\"<>,./?~"`},
		// Seen in the wild:
		// Cloudflare: (verify "$" edge case)
		{`cf1`, []string{`302,test2.foo.com,https://goo.com/$1`}, `"302,test2.foo.com,https://goo.com/$1"`},
		// Cloudflare: (interior quotes)
		{`cf2`, []string{`http.host eq "test2.foo.com" and http.request.uri.path eq "/"`}, `"http.host eq \"test2.foo.com\" and http.request.uri.path eq \"/\""`},
		// Vercel: (URLs in TXT records)
		{"vercel1", []string{`letsencrypt.org; accounturi=https://acme-v01.api.letsencrypt.org/acme/reg/1234567`}, `"letsencrypt.org; accounturi=https://acme-v01.api.letsencrypt.org/acme/reg/1234567"`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rd := rdata.TXT{Txt: tc.txts}

			// Check the primary result:

			escaped := rd.String()
			if escaped != tc.want {
				t.Fatalf("expected %s, got %s", tc.want, escaped)
			}

			// Round-trip:

			plain, err := dns.NewData(dns.TypeTXT, escaped)
			if err != nil {
				t.Fatalf("roundtrip: unexpected error: %s", err)
			}
			rdtxts, ok := plain.(rdata.TXT)
			if !ok {
				t.Fatalf("roundtrip: expected type %T, got %T", rdata.TXT{}, rdtxts)
			}
			if !slices.Equal(rdtxts.Txt, tc.txts) {
				t.Fatalf("roundtrip: expected %s, got %s", tc.txts, rdtxts.Txt)
			}
		})
	}
}
