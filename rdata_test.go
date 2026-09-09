package dns_test

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/dnsutil"
	"codeberg.org/miekg/dns/internal/dnsstring"
	"codeberg.org/miekg/dns/rdata"
)

// Example on how get the text presentation of a [dns.RR].
func ExampleRDATA_string() {
	rr := &dns.MX{Hdr: dns.Header{Name: "miek.nl.", Class: dns.ClassINET, TTL: 3600}, Preference: 10, Mx: "mx.miek.nl."}
	s := rr.Header().String() + " " + dnsutil.TypeToString(dns.RRToType(rr)) + "\t" + rr.Data().String()
	fmt.Println(s)
	// Output: miek.nl.	3600	IN MX	10 mx.miek.nl.
}

// Example on how to set the rdata of an RR.
func ExampleRDATA_newData() {
	rd, _ := dns.NewData(dns.TypeMX, "10 mx.miek.nl.")
	rr := dns.TypeToRR[dns.TypeMX]()
	rr.Header().Name = "miek.nl."
	rr.Header().Class = dns.ClassINET
	fn := dns.TypeToRDATA[dns.TypeMX]
	// Set the rdata in the rr.
	fn(rr, rd)
	fmt.Println(rr)
	// Output: miek.nl.	0	IN	MX	10 mx.miek.nl.
}

func TestTypeToRDATA(t *testing.T) {
	testcases := []struct {
		name string
		t    uint16
		in   string
		fn   func(rr dns.RR) error
	}{
		{
			"mx",
			dns.TypeMX,
			"10 mx.miek.nl.",
			func(rr dns.RR) error {
				mx, ok := rr.(*dns.MX)
				if !ok {
					return fmt.Errorf("expected MX, got %T", rr)
				}
				if mx.Preference != 10 {
					return fmt.Errorf("expected 10, got %d", mx.Preference)
				}
				if mx.Mx != "mx.miek.nl." {
					return fmt.Errorf("expected mx.miek.nl., got %s", mx.Mx)
				}
				return nil
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rd, _ := dns.NewData(tc.t, tc.in, ".")
			rr := dns.TypeToRR[tc.t]()
			fn := dns.TypeToRDATA[tc.t]
			fn(rr, rd)

		})
	}
}

func TestNewData(t *testing.T) {
	testcases := []struct {
		name string
		t    uint16
		in   string
		fn   func(rd dns.RDATA) error
	}{
		{
			"mx-origin-ok",
			dns.TypeMX,
			"10 mx.miek.nl",
			func(rd dns.RDATA) error {
				if rd == nil {
					return fmt.Errorf("expected rd, got none")
				}
				mx := rd.(rdata.MX)
				if mx.Preference != 10 {
					return fmt.Errorf("expected 10, got %d", mx.Preference)
				}
				if mx.Mx != "mx.miek.nl." {
					return fmt.Errorf("expected mx.miek.nl., got %s", mx.Mx)
				}
				return nil
			},
		},
		{
			"mx-ok",
			dns.TypeMX,
			"10 mx.miek.nl.",
			func(rd dns.RDATA) error {
				if rd == nil {
					return fmt.Errorf("expected rd, got none")
				}
				mx := rd.(rdata.MX)
				if mx.Preference != 10 {
					return fmt.Errorf("expected 10, got %d", mx.Preference)
				}
				if mx.Mx != "mx.miek.nl." {
					return fmt.Errorf("expected mx.miek.nl., got %s", mx.Mx)
				}
				return nil
			},
		},
		{
			"mx-space-fail",
			dns.TypeMX,
			" 10 mx.miek.nl.",
			func(rd dns.RDATA) error {
				if rd.(rdata.MX).Preference == 0 {
					return nil
				}
				return fmt.Errorf("expected nil rd: %v", rd)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rd, _ := dns.NewData(tc.t, tc.in, ".")
			if err := tc.fn(rd); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDataDrift(t *testing.T) {
	// The rdata package renders algorithms via its own copy of the map, which must not drift.
	if !maps.Equal(dns.AlgorithmToString, dnsstring.AlgorithmToString) {
		t.Errorf("expected %v, got %v", dns.AlgorithmToString, dnsstring.AlgorithmToString)
	}
}

// TestTXT round-trips TXT records.
func TestTXT(t *testing.T) {
	testcases := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		input   string
		want    []string
		wantErr bool
	}{
		// Basic:
		{"unquoted", `unquoted`, []string{`unquoted`}, false},
		{"quoted", `"quoted"`, []string{`quoted`}, false},
		{"quoted 3", `"one" "two" "three"`, []string{`one`, `two`, `three`}, false},
		// whitespace:
		{"whitespace", `  whitespace  `, []string{`whitespace`}, false},
		{"quoted ws", `" white space "`, []string{` white space `}, false},
		// quotes and semicolons:
		{"dquote", `"double\"quote"`, []string{`double"quote`}, false},
		{"semi", `;`, []string{}, false},
		{"escaped semi", `\;`, []string{`;`}, false},
		{"quoted escaped semi", `"\;"`, []string{`;`}, false},
		{"ddd semi", `\059`, []string{`;`}, false},
		// `m`: (a normal ascii char)
		{"m", `\m`, []string{`m`}, false},
		{"quoted escaped m", `"\m"`, []string{`m`}, false},
		{"ddd m", `\109`, []string{`m`}, false},
		{"quoted ddd m", `"\109"`, []string{`m`}, false},
		// `\003`: (a non-printable ascii char)
		{"etx", `\003`, []string{"\x03"}, false},
		{"quoted etx", `"\003"`, []string{"\x03"}, false},
		// `\127`: (a non-printable ascii char, highest value in the ascii table)
		{"del", `\127`, []string{"\x7f"}, false},
		{"quoted del", `"\127"`, []string{"\x7f"}, false},
		{"swearing", `"!@#$%^&*();:'\"<>,./?~"`, []string{`!@#$%^&*();:'"<>,./?~`}, false},
		//
		// Special cases
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
		//
		// Errors:
		{"unclosedquote", `"unclosedquote`, []string{`unclosedquote`}, true},
		{"unstartedquote", `unstartedquote"`, []string{`unstartedquote"`}, true},
		{"unfinished slash", `\`, []string{`\`}, true},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			// Check the primary result:
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
			rdtxts, ok := rd.(rdata.TXT)
			if !ok {
				t.Fatalf("expected type %T, got %T", rdata.TXT{}, rdtxts)
			}
			if !slices.Equal(rdtxts.Txt, tc.want) {
				t.Fatalf("expected %s, got %s", tc.want, rdtxts.Txt)
			}

			// Round-trip:
			escaped := rdtxts.String()
			rd2, err := dns.NewData(dns.TypeTXT, escaped)
			if err != nil {
				t.Fatalf("roundtrip: unexpected error: %s", err)
			}
			rdtxts2, ok := rd2.(rdata.TXT)
			if !ok {
				t.Fatalf("roundtrip: expected type %T, got %T", rdata.TXT{}, rdtxts2)
			}
			if !slices.Equal(rdtxts.Txt, rdtxts2.Txt) {
				t.Errorf("roundtrip: expected %s, got %s", rdtxts.Txt, rdtxts2.Txt)
			}

			// Round-trip string representation:
			rt := rd2.String()
			if escaped != rt {
				t.Errorf("2nd round trip: expected %s, got %s", escaped, rt)
			}

		})
	}
}
