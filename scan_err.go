package dns

func (zp *ZoneParser) ResetErr() ([]string, error) {
	tokens := make([]string, 0)
	if zp.parseErr != nil {
		nextLine := false
		for !nextLine {
			l, _ := zp.c.Next()
			tokens = append(tokens, l.token)
			if l.value == zNewline || l.value == zEOF {
				nextLine = true
			}
		}
		zp.parseErr = nil
	}
	return tokens, nil
}
