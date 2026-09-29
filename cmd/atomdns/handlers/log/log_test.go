package log_test

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/lanrat/dns/cmd/atomdns/atomtest"
	"github.com/lanrat/dns/cmd/atomdns/handlers/log"
	"github.com/lanrat/dns/cmd/atomdns/internal/dnsctx"
	"github.com/lanrat/dns/dnstest"
)

func TestLog(t *testing.T) {
	l := &log.Log{
		Contexts: map[string][]string{
			"hello": {"here", "there"},
		},
	}

	b := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(b, nil))
	slog.SetDefault(logger)
	m := dnstest.NewMsg()

	ctx := context.Background()
	ctx = dnsctx.WithValue(ctx, "hello/here", "not far")
	ctx = dnsctx.WithValue(ctx, "hello/there", "far")

	testcases := []struct {
		name string
		exp  []string
	}{
		{"ipv4",
			[]string{`hello.here="not far" hello.there=far`, "remote=198.51.100.1"},
		},
		{"ipv6",
			[]string{`hello.here="not far" hello.there=far`, "remote=2001:db8::1"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			b.Reset()
			tw := dnstest.NewTestRecorder()
			if tc.name == "ipv6" {
				tw = dnstest.NewTestRecorder6()
			}
			l.HandlerFunc(atomtest.Echo).ServeDNS(ctx, tw, m)
			s := b.String()
			for _, exp := range tc.exp {
				if !strings.Contains(s, exp) {
					t.Fatalf("expected %s, got none", exp)
				}
			}
		})
	}
}
