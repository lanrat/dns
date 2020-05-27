[![Go Report Card](https://goreportcard.com/badge/github.com/lanrat/dns)](https://goreportcard.com/report/lanrat/dns)
[![Go Doc](https://godoc.org/github.com/lanrat/dns?status.svg)](https://pkg.go.dev/github.com/lanrat/dns?tab=doc)

# A more forgiving fork of [miekg/dns](https://github.com/miekg/dns)

This is a fork of [miekg's wonderful dns library](https://github.com/miekg/dns) `github.com/miekg/dns`

The `miekg/dns` library's policy is "garbage in, garbage out", which works well when your inputs are well sanitized and under control. Unfortunately some real-world use cases supply data in a non-ideal format.

The [miekg branch](https://github.com/lanrat/dns/tree/miekg) of this repository will be kept in sync with the [upstream master](https://github.com/miekg/dns/tree/master) and used to keep this fork's master branch up to date.

Individual feature branches of this repository are used to test single features that modify `miekg/dns` in a single way and one way only. After they have been proven stable they may be merged into master. The goal is to have this fork pass all tests from `miekg/dns` and any additional tests that are added, and be fully backwards compatible with `miekg/dns`.

Some features may be sent as a pull request upstream. If they are accepted then they are removed from the differences listed below.

# Differences from miekg/dns

 * `ResetErr()` added to the ZoneParser, allowing the parser to recover from errors
 * `TypeBitMap` is automatically sorted when packing `CSYNC`, `NSEC`, and `NSEC3` types
 * `AAAA.String()` method converts IPv4 addresses to IPv6 notation
 * `SetDomainFunc()` added to add a function to clean or sanitize zone/domain names as they are parsed
