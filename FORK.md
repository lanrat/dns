# Fork of codeberg.org/miekg/dns

This is a fork of [miekg's dns library](https://codeberg.org/miekg/dns) (v2) with a small set of
changes needed to parse real-world zone files, which are often not clean. It replaces the v1 fork
[github.com/lanrat/dns](https://github.com/lanrat/dns).

The upstream `README.md` is left untouched so merges from upstream never conflict. Fork changes are
kept in new files where possible for the same reason.

## Differences from upstream

- `ZoneParser.ResetErr()` clears a parse error and discards the rest of the bad RR, so parsing can
  continue with the next record instead of stopping at the first error.
  Unlike the v1 fork, it does not discard anything when the error token already ended the RR. The v1
  version dropped the next good record in that case (for example after `b IN A` with no rdata).
  Files: `scan_reset.go`, `scan_reset_test.go`.

## Using the fork

The module path is not renamed. It stays `codeberg.org/miekg/dns`, which keeps the diff against
upstream to the files above. Code imports `codeberg.org/miekg/dns` as normal, and the application's
`go.mod` points that path at this fork:

```text
require codeberg.org/miekg/dns v0.6.117

replace codeberg.org/miekg/dns => <this fork's module location> <version>
```

A `replace` directive only applies in the main module, so this works for applications but not for
libraries imported by other modules.

## Branches

- `miekg`: a mirror of upstream `main`. Never commit to it.
- One feature branch per change, containing only that change (`parseRecover`).
- `main`: `miekg` with every feature branch merged in, plus this file.

## Updating from upstream

```shell
# setup, once
git remote add upstream https://codeberg.org/miekg/dns.git

# update the mirror branch
git fetch upstream
git checkout miekg
git merge --ff-only upstream/main
git push origin miekg

# merge into main and verify
git checkout main
git merge miekg
go test ./...
git push origin main
```

To change a feature, commit on its branch and merge the branch into `main` again.

## Features of the v1 fork that were not ported

These existed in the v1 fork and were left out because the application no longer needs them. The
notes are here so they can be re-added if that changes.

### TypeBitMap sorting

The v1 fork sorted the `TypeBitMap` of NSEC, NSEC3 and CSYNC records while parsing. zonetools sorts
these itself during normalization (`parser/normalize.go`), so the library change was redundant.

### SetDomainFunc (domain name hook)

The v1 fork had `SetDomainFunc(func(string) string)`, a package-level hook applied to the zone origin
and to every domain name before it was validated. zonetools used it to lowercase names and convert
UTF-8 names to punycode (`parser.CleanDomain`). It was dropped for v2 and zonetools now cleans names
after parsing instead.

**What the hook still did that post-parse cleaning cannot.** A name is validated before cleaning, so
a raw UTF-8 label longer than 63 bytes is rejected even if its punycode form would fit. For example,
the Thai label `ที่ปรึกษาธุรกิจ-แผนการตลาด` is 76 bytes in UTF-8 but valid as punycode. Any fully
qualified name with such a label is rejected, as an owner name or in rdata, and the record is lost
with an error such as:

```text
dns: bad owner name: "ที..." at line: 1:59
dns: bad NS Ns: "ns.ที..." at line: 3:84
```

A relative name whose last label is too long does parse, because upstream `IsName` does not check
the last label of a relative name before the origin is appended.

**Why it was dropped.** Labels over 63 octets are invalid DNS (RFC 1035), and IDNs are published as
A-labels (RFC 5890). Zones collected by AXFR or walking come from the wire, so they cannot contain
such labels. Only text exports from registries (CZDS, FTP) could, and only if the export is buggy. In
2026-09, a scan of the zone file archive found no raw non-ASCII names: every file collected on
2026-09-25, all of 2020 to 2023, and one day each from 2024 and 2025.

**How to re-add it.** Two hook points cover every name the zone parser reads:

1. `NewZoneParser` in `scan.go`, applied to `origin` before `dnsutilFqdn`.
2. `Absolute` in `dnsutil/shared.go`, applied to the name before it is validated. This function is
   copied into `zdnsutil.go` in the root, `rdata`, `svcb` and `deleg` packages by
   `dnsutil_generate.go`, so edit `dnsutil/shared.go` and run `go generate`. Do not edit the
   generated copies.

Put the hook variable and setter in a new file so upstream merges stay conflict-free.
