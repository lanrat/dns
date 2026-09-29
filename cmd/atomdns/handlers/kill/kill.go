package kill

import (
	"github.com/lanrat/dns"
)

type Kill int

func (k *Kill) HandlerFunc(next dns.HandlerFunc) dns.HandlerFunc { return nil }
