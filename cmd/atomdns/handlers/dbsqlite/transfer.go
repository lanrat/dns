package dbsqlite

import (
	"context"
	"io"
	"log/slog"
	"path/filepath"

	"github.com/lanrat/dns"
	"github.com/lanrat/dns/cmd/atomdns/internal/dnslog"
	"github.com/lanrat/dns/cmd/atomdns/internal/dnszone"
	"github.com/lanrat/dns/dnsutil"
)

func (d *Dbsqlite) HandlerFuncTransfer(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) {
	z := d.Zones[dns.Zone(ctx)]

	if d.To == nil || z == nil || dnsutil.Network(w) == "udp" {
		m := new(dns.Msg)
		dnsutil.SetReply(m, r)
		m.Rcode = dns.RcodeRefused
		m.Data = r.Data

		if err := m.Pack(); err != nil {
			dnslog.PackFail(ctx, log(), Err(err))
			return
		}
		io.Copy(w, m)
		return
	}
	if err := dnszone.TransferOut(z, ctx, w, r); err != nil {
		log().Debug("Failure to transfer out", Err(err))
		return
	}
	alog := log().With(slog.String("zone", z.Origin()), slog.String("path", filepath.Base(d.Path)), slog.Any("upstream", w.RemoteAddr()), slog.Uint64("serial", uint64(dnszone.Serial(z))))
	alog.Info("Successful transfer out")
}
