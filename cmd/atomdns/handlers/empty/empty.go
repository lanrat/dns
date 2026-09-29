package empty

import (
	"github.com/lanrat/dns"
)

type Empty int

func (e *Empty) HandlerFunc(next dns.HandlerFunc) dns.HandlerFunc { return nil }
