package tsig

import (
	"encoding/base64"
	"fmt"

	"github.com/lanrat/dns"
	"github.com/lanrat/dns/cmd/atomdns/internal/dnsserver"
	"github.com/lanrat/dns/dnsutil"
)

func (t *Tsig) Setup(co *dnsserver.Controller) error {
	if co.Next() {
		args := co.RemainingArgs()
		if len(args) != 3 {
			return co.ArgErr()
		}
		name := dnsutil.Canonical(args[0])
		algo := dnsutil.Canonical(args[1])
		if !dnsutil.IsName(name) {
			return co.PropErr(fmt.Errorf("name %s is not a domain name", name))
		}
		if !dnsutil.IsName(algo) {
			return co.PropErr(fmt.Errorf("algorithm %s is not a domain name", algo))
		}
		if _, err := base64.StdEncoding.DecodeString(args[2]); err != nil {
			return co.PropErr(err)
		}
		t.TSIG = dns.NewTSIG(name, algo, 0)
		t.TSIGSecret = args[2]
	}
	return nil
}
