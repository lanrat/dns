# dns

## A more forgiving fork of [miekg/dns](https://codeberg.org/miekg/dns) (v2)

This is a fork of [miekg's wonderful dns library](https://codeberg.org/miekg/dns) `codeberg.org/miekg/dns`, the v2 rewrite of `github.com/miekg/dns`.

The `miekg/dns` library's policy is "garbage in, garbage out", which works well when your inputs are well sanitized and under control. Unfortunately some real-world use cases supply data in a non-ideal format, such as zone files published by registries.

The [v2/miekg branch](https://github.com/lanrat/dns/tree/v2/miekg) of this repository is kept in sync with the [upstream main](https://codeberg.org/miekg/dns/src/branch/main) and used to keep this fork's `main` branch up to date.

Individual feature branches of this repository (`v2/parseRecover`, `v2/caseInsensitive`, `v2/domainFunc`) are used to test single features that modify `miekg/dns` in a single way and one way only. After they have been proven stable they are merged into `main`. The goal is to have this fork pass all tests from `miekg/dns` and any additional tests that are added, and be fully backwards compatible with `miekg/dns`.

Some features may be sent as a pull request upstream. If they are accepted then they are removed from the differences listed below.

The earlier fork of the v1 library (`github.com/miekg/dns`) is on the [v1 branch](https://github.com/lanrat/dns/tree/v1). Its support branches keep their original names, which is why the v2 branches are prefixed with `v2/`.

### Differences from miekg/dns

* `ResetErr()` added to the ZoneParser, allowing the parser to recover from errors
* Classes and type mnemonics are accepted in any case, as RFC 1035 requires
* `SetDomainFunc()` added to add a function to clean or sanitize zone/domain names as they are parsed
* The module path is `github.com/lanrat/dns` (see [Using this fork](#using-this-fork))

#### `ResetErr()`

Clears a parse error and discards the rest of the bad RR, so parsing continues with the next record instead of stopping at the first error. Unlike the v1 fork, it does not discard anything when the error token already ended the RR; the v1 version dropped the next good record in that case (for example after `b IN A` with no rdata).

Files: `scan_reset.go`, `scan_reset_test.go`.

#### Case-insensitive classes and types

RFC 1035 Section 5.1 makes class and type mnemonics case-insensitive, and some published zone files are entirely lower case (`a. 3600 in soa ...`). Upstream looks up types case-insensitively, but not classes, the `TYPExxx` and `CLASSxxx` prefixes, RRSIG covered types, NSEC/NSEC3/CSYNC type bitmaps or the DSYNC type, so every line of such a file failed to parse.

Files: `internal/dnslex/fold.go`, `scan_case_test.go`, and small edits to `internal/dnslex/lex.go` and `scan_rdata.go`.

#### `SetDomainFunc()`

Sets a hook the zone parser applies to the origin and to every name it reads, before the name is validated. Because it runs before validation, a name that is only valid after the hook's conversion is still accepted: for example the Thai label `ที่ปรึกษาธุรกิจ-แผนการตลาด` is 76 bytes in UTF-8, too long for a DNS label, but 42 as the punycode `xn----twfab7a0egkknqr9fcg7a9c1jid6anz6c8o7f`, a valid name that can be registered and resolved. Without the hook such names fail with `dns: bad owner name` or `dns: bad NS Ns`.

The hook runs in `Absolute` in `dnsutil/shared.go`, which `dnsutil_generate.go` copies into `zdnsutil.go` in the root, `deleg`, `svcb` and `rdata` packages, and in `NewZoneParser` in `scan.go` for the origin. Only the root package's copy is set, so the public `dnsutil.Absolute` is unaffected.

Files: `domainfunc.go`, `domainfunc_test.go`, `domainfunc_scope_test.go`, 3 lines in `dnsutil/shared.go` plus its generated copies, and 3 lines in `scan.go`.

### Not carried over from the v1 fork

* `TypeBitMap` sorting when packing `CSYNC`, `NSEC`, and `NSEC3` types. zonetools sorts these itself during normalization, so the library change was redundant.

## Using this fork

Like the v1 fork, the module path is renamed to `github.com/lanrat/dns`, so it is used like any other module:

```shell
go get github.com/lanrat/dns@main
```

```go
import "github.com/lanrat/dns"
```

The rename is done by `rename-module.sh`, which rewrites `codeberg.org/miekg/dns` to `github.com/lanrat/dns` in `go.mod` and every Go file. The `v2/miekg` mirror and the feature branches keep the upstream path, so the script is run again after merging any of them into `main`.

The `v1` branch has the same module path with the v1 API. Applications pick one by version: a pseudo-version of a `main` commit is the v2 fork, one of a `v1` commit is the v1 fork.

## Updating from upstream

```shell
# setup
git remote add upstream https://codeberg.org/miekg/dns.git
git pull origin main
git pull origin v2/miekg

# update v2/miekg branch from upstream
git checkout v2/miekg
git pull upstream main
git push origin v2/miekg

# merge into main
git checkout main
git merge v2/miekg
git checkout --ours README.md
git add README.md
# fix any merge conflicts; for conflicts only in import paths, take upstream's side
# then rename the module path again in anything the merge brought in
./rename-module.sh
git add .
go test ./...
git commit
git push origin main
```

If the merge conflicts in a generated `zdnsutil.go`, take upstream's version of the generated files, keep this fork's lines in `dnsutil/shared.go`, and run `go run dnsutil_generate.go` to regenerate them.

To change a feature, commit on its branch, merge the branch into `main` again, and run `./rename-module.sh`.
