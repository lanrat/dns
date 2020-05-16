package dns

import "log"

func (zp *ZoneParser) ResetErr() error {
	if zp.parseErr != nil {
		nextLine := false
		tokens := make([]string, 0)
		for !nextLine {
			l, _ := zp.c.Next()
			//log.Printf("miekg reset slurp: %+v\n", l.token)
			tokens = append(tokens, l.token)
			if l.value == zNewline || l.value == zEOF {
				nextLine = true
			}
		}
		zp.parseErr = nil
		// TODO return the line?
		log.Printf("Skipped line: %v", tokens)
	}
	return nil
}
