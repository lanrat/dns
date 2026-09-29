#!/bin/sh
# rename-module.sh rewrites the upstream module path, codeberg.org/miekg/dns, to this fork's,
# github.com/lanrat/dns, in go.mod and every Go file (imports, generator templates and doc links).
# Run it after merging v2/miekg or a feature branch into main, since those still use the upstream
# path. It is idempotent.
set -eu
cd "$(dirname "$0")"
grep -rlZ --include='*.go' --include='go.mod' 'codeberg.org/miekg/dns' . |
	xargs -0 -r sed -i 's#codeberg\.org/miekg/dns#github.com/lanrat/dns#g'
gofmt -l -w . >/dev/null
