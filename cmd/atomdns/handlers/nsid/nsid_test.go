package nsid_test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/lanrat/dns"
	"github.com/lanrat/dns/cmd/atomdns/atomtest"
	"github.com/lanrat/dns/cmd/atomdns/handlers/nsid"
	"github.com/lanrat/dns/dnstest"
)

func TestNsid(t *testing.T) {
	in := "Use the force"
	n := &nsid.Nsid{Data: hex.EncodeToString([]byte(in))}

	m := dnstest.NewMsg()
	m.Pseudo = []dns.RR{&dns.NSID{}}
	m.Pack()

	tw := dnstest.NewTestRecorder()
	n.HandlerFunc(atomtest.Echo).ServeDNS(context.TODO(), tw, m)

	tw.Msg.Unpack()
	if len(tw.Msg.Pseudo) != 1 {
		t.Fatal("expected pseudo section")
	}
	if tw.Msg.Pseudo[0].(*dns.NSID).Nsid != hex.EncodeToString([]byte(in)) {
		t.Fatalf("expected NSID RR contain: %s", in)
	}
}
