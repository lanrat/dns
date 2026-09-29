package empty

import (
	"github.com/lanrat/dns/cmd/atomdns/internal/dnsserver"
)

func (e *Empty) Setup(co *dnsserver.Controller) error {
	for co.Next() {
	}
	return nil
}
