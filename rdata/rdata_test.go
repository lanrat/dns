package rdata_test

import (
	"slices"
	"strings"
	"testing"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/rdata"
)

// TestTXT round-trips TXT records.
func TestTXT(t *testing.T) {
	testcases := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{"unquoted", `unquoted`, []string{`unquoted`}, false},
		{"quoted", `"quoted"`, []string{`quoted`}, false},
		{"quoted", `"one" "two" "three"`, []string{`one`, `two`, `three`}, false},

		{"whitespace", `  whitespace  `, []string{`whitespace`}, false},
		{"quoted ws", `" white space "`, []string{` white space `}, false},

		{"dquote", `"double\"quote"`, []string{`double"quote`}, false},
		{"semi", `;`, []string{}, false},
		{"escaped semi", `\;`, []string{`;`}, false},
		{"quoted escaped semi", `"\;"`, []string{`;`}, false},
		{"ddd semi", `\059`, []string{`;`}, false},

		{"m", `\m`, []string{`m`}, false},
		{"quoted escaped m", `"\m"`, []string{`m`}, false},
		{"ddd m", `\109`, []string{`m`}, false},
		{"quoted ddd m", `"\109"`, []string{`m`}, false},

		{"etx", `\003`, []string{"\x03"}, false},
		{"quoted etx", `"\003"`, []string{"\x03"}, false},

		{"del", `\127`, []string{"\x7f"}, false},
		{"quoted del", `"\127"`, []string{"\x7f"}, false},

		// long strings
		{"long", `"` + strings.Repeat("x", 255) + `"`, []string{strings.Repeat("x", 255)}, false},
		{"too long", `"` + strings.Repeat("x", 255) + strings.Repeat("y", 10) + `"`,
			[]string{strings.Repeat("x", 255), strings.Repeat("y", 10)}, false},
		//
		// AWS Route53 is known to produce TXT records with no space between the segments if they are both quoted.
		{"amazon", `"first""second"`, []string{`first`, `second`}, false},
		// Vercel has unescaped and escaped semicolons.
		{"vercel1", `"letsencrypt.org; accounturi=https://acme-v01.api.letsencrypt.org/acme/reg/1234567"`,
			[]string{`letsencrypt.org; accounturi=https://acme-v01.api.letsencrypt.org/acme/reg/1234567`}, false},
		{"vercel2", `"letsencrypt.org\; accounturi=https://acme-v01.api.letsencrypt.org/acme/reg/1234567"`,
			[]string{`letsencrypt.org; accounturi=https://acme-v01.api.letsencrypt.org/acme/reg/1234567`}, false},
		// Cloudflare has interior quotes and a $1 in the string.
		{`cf url`, `"302,test2.foo.com,https://goo.com/$1"`,
			[]string{`302,test2.foo.com,https://goo.com/$1`}, false},
		{`cf quotes`, `"http.host eq \"test2.foo.com\" and http.request.uri.path eq \"/\""`,
			[]string{`http.host eq "test2.foo.com" and http.request.uri.path eq "/"`}, false},

		{"unclosedquote", `"unclosedquote`, []string{`unclosedquote`}, true},
		{"unstartedquote", `unstartedquote"`, []string{`unstartedquote"`}, true},
		{"unfinished slash", `\`, []string{`\`}, true},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			rd, err := dns.NewData(dns.TypeTXT, tc.input)
			if tc.wantErr == false && err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if tc.wantErr == true && err == nil {
				t.Fatalf("expected error, got none")
			}
			if err != nil {
				return
			}
			rdtxts, _ := rd.(rdata.TXT)
			if !slices.Equal(rdtxts.Txt, tc.want) {
				t.Fatalf("expected %s, got %s", tc.want, rdtxts.Txt)
			}

			escaped := rdtxts.String()
			rd2, err := dns.NewData(dns.TypeTXT, escaped)
			if err != nil {
				t.Fatalf("roundtrip: unexpected error: %s", err)
			}
			rdtxts2, _ := rd2.(rdata.TXT)
			if !slices.Equal(rdtxts.Txt, rdtxts2.Txt) {
				t.Errorf("roundtrip: expected %s, got %s", rdtxts.Txt, rdtxts2.Txt)
			}

			// Round-trip string representation:
			if rt := rd2.String(); escaped != rt {
				t.Errorf("2nd round trip: expected %s, got %s", escaped, rt)
			}
		})
	}
}
