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
- Classes and type mnemonics are accepted in any case, as RFC 1035 Section 5.1 requires. Upstream
  looks up types case-insensitively but not classes, the `TYPExxx` and `CLASSxxx` prefixes, RRSIG
  covered types, NSEC/NSEC3/CSYNC type bitmaps or the DSYNC type. Some published zone files are
  entirely lower case (`a. 3600 in soa ...`), and every line of those failed to parse. This is a bug
  fix worth sending upstream, after which this change can be dropped.
  Files: `internal/dnslex/fold.go`, `scan_case_test.go`, and small edits to `internal/dnslex/lex.go`
  and `scan_rdata.go`.
- `SetDomainFunc(func(string) string)` sets a hook the zone parser applies to the origin and to every
  name it reads, before the name is validated. zonetools uses it to lowercase names and convert UTF-8
  names to punycode. Because it runs before validation, a UTF-8 label longer than 63 bytes whose
  punycode form fits is accepted: for example the Thai label `ที่ปรึกษาธุรกิจ-แผนการตลาด` is 76 bytes in
  UTF-8 and 42 as `xn----twfab7a0egkknqr9fcg7a9c1jid6anz6c8o7f`, a valid name that can be registered
  and resolved. Without the hook such names fail with `dns: bad owner name` or `dns: bad NS Ns`.
  The hook runs in `Absolute` in `dnsutil/shared.go`, which `dnsutil_generate.go` copies into
  `zdnsutil.go` in the root, `deleg`, `svcb` and `rdata` packages, and in `NewZoneParser` in
  `scan.go` for the origin. Only the root package's copy is set, so the public `dnsutil.Absolute` is
  unaffected. Files: `domainfunc.go`, `domainfunc_test.go`, `domainfunc_scope_test.go`, 3 lines in
  `dnsutil/shared.go` plus its generated copies, and 3 lines in `scan.go`.

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
- One feature branch per change, containing only that change (`parseRecover`, `caseInsensitive`,
  `domainFunc`).
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

If the merge conflicts in a generated `zdnsutil.go`, take upstream's version of the generated files,
keep the fork's lines in `dnsutil/shared.go`, and run `go run dnsutil_generate.go` to regenerate
them.

To change a feature, commit on its branch and merge the branch into `main` again.

## Feature of the v1 fork that was not ported

This existed in the v1 fork and was left out because the application no longer needs it. The note is
here so it can be re-added if that changes.

### TypeBitMap sorting

The v1 fork sorted the `TypeBitMap` of NSEC, NSEC3 and CSYNC records while parsing. zonetools sorts
these itself during normalization (`parser/normalize.go`), so the library change was redundant.
