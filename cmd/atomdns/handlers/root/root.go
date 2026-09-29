package root

import "github.com/lanrat/dns"

type Root struct {
	global string
}

func (r *Root) HandlerFunc(next dns.HandlerFunc) dns.HandlerFunc { return nil }
