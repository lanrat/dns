package cookie_test

import (
	"context"
	"encoding/hex"
	"hash/fnv"
	"io"
	"testing"

	"github.com/lanrat/dns"
	"github.com/lanrat/dns/cmd/atomdns/atomtest"
	"github.com/lanrat/dns/cmd/atomdns/handlers/cookie"
	"github.com/lanrat/dns/dnstest"
)

func TestCookie(t *testing.T) {
	c := &cookie.Cookie{Secret: "geheim"}

	f := fnv.New64()
	io.WriteString(f, "::1")
	io.WriteString(f, "::1")
	io.WriteString(f, "ook geheim")
	cookie := &dns.COOKIE{Cookie: hex.EncodeToString(f.Sum(nil))}

	m := dnstest.NewMsg()
	m.Pseudo = []dns.RR{cookie}
	m.Pack()

	tw := dnstest.NewTestRecorder()
	c.HandlerFunc(atomtest.Echo).ServeDNS(context.TODO(), tw, m)

	tw.Msg.Unpack()
	if len(tw.Msg.Pseudo) != 1 {
		t.Fatal("expected pseudo section")
	}
	if _, ok := tw.Msg.Pseudo[0].(*dns.COOKIE); !ok {
		t.Fatal("expected COOKIE RR")
	}
}
