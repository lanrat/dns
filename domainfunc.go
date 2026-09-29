package dns

// SetDomainFunc sets a function the zone parser applies to every domain name it reads, including the
// origin, before the name is validated. Use it to normalize names, for example lower-casing them or
// converting UTF-8 names to punycode, so that a name which is only valid after the conversion (a
// UTF-8 label longer than 63 bytes whose punycode form fits) is still accepted.
//
// The function is global to the package and not safe to change while zone files are being parsed.
// Pass nil to remove it.
func SetDomainFunc(f func(string) string) { domainFunc = f }
