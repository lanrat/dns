package dns

import (
	"slices"
	"strings"
	"testing"
)

func TestZoneParserResetErr(t *testing.T) {
	tests := []struct {
		name   string
		zone   string
		owners []string
		errs   int
	}{
		{
			name:   "no errors",
			zone:   "a IN A 127.0.0.1\nb IN A 127.0.0.2\n",
			owners: []string{"a.example.", "b.example."},
		},
		{
			name:   "bad rdata mid line",
			zone:   "a IN A 127.0.0.1\nb IN A bad 127.0.0.2\nc IN A 127.0.0.3\n",
			owners: []string{"a.example.", "c.example."},
			errs:   1,
		},
		{
			name:   "garbage after rdata",
			zone:   "a IN A 127.0.0.1\nb IN A 127.0.0.2 extra junk\nc IN A 127.0.0.3\n",
			owners: []string{"a.example.", "c.example."},
			errs:   1,
		},
		{
			name:   "unknown type",
			zone:   "a IN A 127.0.0.1\nb IN FOOBAR x y z\nc IN A 127.0.0.3\n",
			owners: []string{"a.example.", "c.example."},
			errs:   1,
		},
		{
			name:   "missing rdata",
			zone:   "a IN A 127.0.0.1\nb IN A\nc IN A 127.0.0.3\n",
			owners: []string{"a.example.", "c.example."},
			errs:   1,
		},
		{
			name:   "bad multi line record",
			zone:   "a IN A 127.0.0.1\nb IN SOA ns. mbox. (\n bad\n 2\n 3\n 4\n 5 )\nc IN A 127.0.0.3\n",
			owners: []string{"a.example.", "c.example."},
			errs:   1,
		},
		{
			name:   "consecutive bad lines",
			zone:   "a IN A bad\nb IN A bad\nc IN A 127.0.0.3\n",
			owners: []string{"c.example."},
			errs:   2,
		},
		{
			name:   "bad last line",
			zone:   "a IN A 127.0.0.1\nb IN A bad",
			owners: []string{"a.example."},
			errs:   1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			zp := NewZoneParser(strings.NewReader(tc.zone), "example.", "")
			owners := []string{}
			errs := 0
			for rr, ok := zp.Next(); ; rr, ok = zp.Next() {
				if !ok {
					if zp.Err() == nil {
						break
					}
					errs++
					if errs > 10 {
						t.Fatalf("too many errors, last: %v", zp.Err())
					}
					zp.ResetErr()
					continue
				}
				owners = append(owners, rr.Header().Name)
			}
			if !slices.Equal(owners, tc.owners) {
				t.Errorf("got owners %v, want %v", owners, tc.owners)
			}
			if errs != tc.errs {
				t.Errorf("got %d errors, want %d", errs, tc.errs)
			}
		})
	}
}

func TestZoneParserResetErrNoError(t *testing.T) {
	zp := NewZoneParser(strings.NewReader("a IN A 127.0.0.1\n"), "example.", "")
	if tokens := zp.ResetErr(); tokens != nil {
		t.Fatalf("expected nil tokens without an error, got %v", tokens)
	}
	if _, ok := zp.Next(); !ok {
		t.Fatalf("expected an RR after ResetErr without an error: %v", zp.Err())
	}
}
