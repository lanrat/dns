package dns

import "github.com/lanrat/dns/internal/dnslex"

// ResetErr clears the current parse error so that parsing can resume with the next RR. The remaining
// tokens of the RR that caused the error are discarded and returned. When the token that caused the
// error ended the RR, nothing further is discarded. ResetErr does nothing when there is no error.
//
// This is used to skip over bad records instead of stopping at the first one:
//
//	for rr, ok := zp.Next(); ; rr, ok = zp.Next() {
//		if !ok {
//			if zp.Err() == nil {
//				break // EOF
//			}
//			log.Println(zp.Err())
//			zp.ResetErr()
//			continue
//		}
//		// Do something with rr
//	}
func (zp *ZoneParser) ResetErr() []string {
	if zp.parseErr == nil {
		if zp.sub != nil {
			return zp.sub.ResetErr()
		}
		return nil
	}

	l := zp.parseErr.lex
	zp.parseErr = nil
	tokens := []string{}
	if l.Value == dnslex.Newline || l.Value == dnslex.EOF {
		return tokens
	}
	for {
		l, _ = zp.c.Next()
		switch l.Value {
		case dnslex.Newline, dnslex.EOF:
			return tokens
		case dnslex.Blank:
		default:
			tokens = append(tokens, l.Token)
		}
	}
}
