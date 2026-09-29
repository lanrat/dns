package ecs_test

import (
	"context"
	"testing"

	"github.com/lanrat/dns"
	"github.com/lanrat/dns/cmd/atomdns/handlers/ecs"
	"github.com/lanrat/dns/cmd/atomdns/internal/dnsctx"
	"github.com/lanrat/dns/dnstest"
	"github.com/lanrat/dns/dnsutil"
)

func TestEcs(t *testing.T) {
	e := &ecs.Ecs{}

	ecs := &dns.SUBNET{Family: dnsutil.IPv4Family, Netmask: 32, Address: dnstest.IPv4}
	m := dnstest.NewMsg()
	m.Pseudo = []dns.RR{ecs}
	m.Pack()

	tw := dnstest.NewTestRecorder()
	next := dns.HandlerFunc(func(ctx context.Context, _ dns.ResponseWriter, _ *dns.Msg) {
		address := dnsctx.Value(ctx, e.Key()+"/addr")
		if address == nil {
			t.Fatal("expected ecs/addr, got none")
		}
		if address != dnstest.IPv4 {
			t.Fatalf("expected %s, got %s", dnstest.IPv4, address)
		}
	})
	ctx := context.TODO()
	e.HandlerFunc(next).ServeDNS(ctx, tw, m)
}
