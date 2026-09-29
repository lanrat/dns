package tsig

import (
	"testing"

	"github.com/lanrat/dns"
	"github.com/lanrat/dns/cmd/atomdns/internal/dnsserver"
)

func TestSetup(t *testing.T) {
	testcases := []struct {
		input string
		exp   *dns.TSIG
		scr   string
	}{
		{
			`tsig example.org hmac-sha512 NoTCJU+DMqFWywaPyxSijrDEA/eC3nK0xi3AMEZuPVk=`,
			dns.NewTSIG("example.org.", "hmac-sha512.", 0), "NoTCJU+DMqFWywaPyxSijrDEA/eC3nK0xi3AMEZuPVk=",
		},
	}
	for i, tc := range testcases {
		tsig := new(Tsig)
		co := dnsserver.NewTestController(tc.input)
		tsig.Setup(co)

		if tc.scr != tsig.TSIGSecret {
			t.Errorf("test %d: expected %s, got %s", i, tc.scr, tsig.TSIGSecret)
		}

		if !dns.Equal(tc.exp, tsig.TSIG) {
			t.Errorf("test %d: expected %s, got %s", i, tc.exp, tsig.TSIG)
		}
	}
}
