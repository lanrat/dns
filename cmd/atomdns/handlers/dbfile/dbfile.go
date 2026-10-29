package dbfile

import (
	"context"
	"io"
	"sync"

	"github.com/lanrat/dns"
	"github.com/lanrat/dns/cmd/atomdns/handlers/dbfile/zone"
	"github.com/lanrat/dns/cmd/atomdns/internal/dnsctx"
	"github.com/lanrat/dns/cmd/atomdns/internal/dnslog"
	"github.com/lanrat/dns/cmd/atomdns/internal/dnszone"
	"github.com/lanrat/dns/dnsutil"
)

type Dbfile struct {
	Path string

	// Zones holds all the zone this instance of Dbfile is called for.
	Zones        map[string]*zone.Zone
	sync.RWMutex // protects Zones

	ctx    context.Context
	cancel context.CancelFunc

	To   *Transfer
	From *Transfer
}

func (d *Dbfile) HandlerFunc(next dns.HandlerFunc) dns.HandlerFunc {
	return dns.HandlerFunc(func(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) {
		if r.Opcode == dns.OpcodeNotify {
			d.HandlerFuncNotify(ctx, w, r)
			return
		}
		if _, qtype := dnsutil.Question(r); qtype == dns.TypeAXFR || qtype == dns.TypeIXFR {
			d.HandlerFuncTransfer(ctx, w, r)
			return
		}

		z := d.Zone(dns.Zone(ctx))
		if z == nil {
			next.ServeDNS(ctx, w, r) // call refused guard
			return
		}

		m := dnszone.Retrieve(z, r, nil)
		if m == nil { // call refused guard
			next.ServeDNS(ctx, w, r)
			return
		}

		m.Data = r.Data
		m = dnsctx.Funcs(ctx, m)
		if err := m.Pack(); err != nil {
			dnslog.PackFail(ctx, log(), Err(err))
			return
		}
		io.Copy(w, m)
	})
}

func (d *Dbfile) Zone(origin string) *zone.Zone {
	d.RLock()
	z := d.Zones[origin]
	d.RUnlock()
	return z
}
