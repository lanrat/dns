package drunk

import (
	"context"
	"io"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/lanrat/dns"
	"github.com/lanrat/dns/cmd/atomdns/internal/dnsctx"
	"github.com/lanrat/dns/cmd/atomdns/internal/dnslog"
	"github.com/lanrat/dns/dnstest"
	"github.com/lanrat/dns/dnsutil"
)

type Drunk struct {
	i        atomic.Uint64 // counter of queries
	drop     uint64
	delay    uint64
	truncate uint64

	duration time.Duration
}

func (d *Drunk) HandlerFunc(next dns.HandlerFunc) dns.HandlerFunc {
	return dns.HandlerFunc(func(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) {
		i := d.i.Load()
		d.i.Add(1)

		drop := d.drop > 0 && i%d.drop == 0
		delay := d.delay > 0 && i%d.delay == 0
		trunc := d.truncate > 0 && i%d.truncate == 0

		m := r.Copy()
		dnsutil.SetReply(m, r)
		m.Authoritative = true
		m.Truncated = trunc

		rw := dnstest.NewRecorder(w)
		next.ServeDNS(ctx, rw, r)

		if drop || rw.Msg == nil { // drop or hijacked conn
			log().With(dnsctx.Id(ctx)).Debug("Dropping")
			return
		}
		if delay {
			log().With(dnsctx.Id(ctx)).Debug("Delaying", slog.Duration("delay", d.duration))
			time.Sleep(d.duration)
		}
		if trunc {
			rw.Msg.Truncated = true
			// have to repack now
			if err := rw.Msg.Pack(); err != nil {
				dnslog.PackFail(ctx, log(), Err(err))
			}
		}

		io.Copy(w, rw.Msg)
	})
}
