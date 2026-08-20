package dns

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/mldsa"
	"crypto/rsa"
	"testing"

	"codeberg.org/miekg/dns/rdata"
)

func TestDNSSECSignVerify(t *testing.T) {
	// Add wildcard, sorting of RRs. etc.
	testcases := []struct {
		name      string
		algorithm uint8
		bitsize   int
		rrs       []RR
	}{
		{
			"rsasha256", RSASHA256, 1024,
			[]RR{
				&SRV{Hdr: Header{Name: "srv.miek.nl.", Class: ClassINET, TTL: 600}, SRV: rdata.SRV{Port: 1000, Weight: 80, Target: "web1.miek.nl."}},
			},
		},
		{
			"ecdsap256sha256", ECDSAP256SHA256, 256,
			[]RR{
				&SRV{Hdr: Header{Name: "srv.miek.nl.", Class: ClassINET, TTL: 600}, SRV: rdata.SRV{Port: 1000, Weight: 80, Target: "web1.miek.nl."}},
			},
		},
		{
			"ed25519", ED25519, 256,
			[]RR{
				&SRV{Hdr: Header{Name: "srv.miek.nl.", Class: ClassINET, TTL: 600}, SRV: rdata.SRV{Port: 1000, Weight: 80, Target: "web1.miek.nl."}},
			},
		},
		{
			"mldsa44", MLDSA44, 256,
			[]RR{
				&SRV{Hdr: Header{Name: "srv.miek.nl.", Class: ClassINET, TTL: 600}, Port: 1000, Weight: 80, Target: "web1.miek.nl."},
			},
		},
		{
			"rsasha256-sorting", RSASHA256, 1024,
			[]RR{
				&NS{Hdr: Header{Name: "miek.nl.", Class: ClassINET, TTL: 600}, NS: rdata.NS{Ns: "linode.atoom.net."}},
				&NS{Hdr: Header{Name: "miek.nl.", Class: ClassINET, TTL: 600}, NS: rdata.NS{Ns: "ns-ext.nlnetlabs.nl."}},
				&NS{Hdr: Header{Name: "miek.nl.", Class: ClassINET, TTL: 600}, NS: rdata.NS{Ns: "omval.tednet.nl"}},
			},
		},
	}

	options := &SignOption{}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var err error

			key := NewDNSKEY("miek.nl.", tc.algorithm)
			priv, _ := key.Generate(tc.bitsize)

			sig := NewRRSIG("miek.nl.", tc.algorithm, key.KeyTag())
			switch tc.algorithm {
			case RSASHA256:
				err = sig.Sign(priv.(*rsa.PrivateKey), tc.rrs, options)
			case ECDSAP256SHA256:
				err = sig.Sign(priv.(*ecdsa.PrivateKey), tc.rrs, options)
			case ED25519:
				err = sig.Sign(priv.(ed25519.PrivateKey), tc.rrs, options)
			case MLDSA44:
				err = sig.Sign(priv.(*mldsa.PrivateKey), tc.rrs, options)
			}
			if err != nil {
				t.Fatalf("failure to sign: %s", err)
			}

			err = sig.Verify(key, tc.rrs, options)
			if err != nil {
				t.Fatalf("failure to verify: %s", err)
			}
		})
	}
}

func TestDNSSECKeyTag(t *testing.T) {
	testcases := []struct {
		tag uint16
		rr  RR
	}{
		{
			10771, dnstestNew("example.net. 3600 IN DNSKEY 257 3 14 xKYaNhWdGOfJ+nPrL8/arkwf2EY3MDJ+SErKivBVSum1w/egsXvSADtNJhyem5RCOpgQ6K8X1DRSEkrbYQ+OB+v8/uX45NBwY8rp65F6Glur8I/mlVNgF6W/qTI37m40"),
		},
		{
			59829, dnstestNew("example.com.	3600 IN	DNSKEY 257 3 18 17K0clSq4NtF55MNSpjSyX2PE5fReJ2voXAksxbpvslPyZRtQvGbeadBO7qjPnFJy0LtURVpOsBB+suYit61/g4dhjEYSZW1ksOX0ilOLhT5CqQUujgmiZrEP0zMrLwm6agyuVEY1ctDPL75ZgsAE44IF/YediyidMNq1VTrIqrBFi5KsBrLoeOMTv2PgLZbMz0PcuVd/nHOnB67mInnxWEGwP1zgDoq7P6v3teqPLLO2lTRK9jNNqeM+XWUO0er0l6ICsRS5XQu0ejRqCr6huWQx1jBWuTShA2SvKGlCQ9ASWWX/KfYuVE/GhvabpUKqpjeRnUH1KT1pPBZkhZYLDVy9i7aiQWrNYFnDEoCd3oz4Mpylf2PT/bRoKOnaD1l9fX3/GDaAj6CbF+SFEwC99G6EHWYdVPqk2f8122ZC3+pnNRa/biDbUPkWfUYffBYR5cJoB6mg1k1+nBGCZDNPcG6QBupS6sd3kGsZ6szGdysoGBI1MTu8n7hOpwX0FOPQw8tZC3CQVZg3niHfY2KvHJSOXjAQuQoX0MZhGxEEmJCl2hEwQ5Va6IVtacZ5Z0MayqW05hZBx/cws3nUkp77a5U6FsxjoVOj+Ky8+36yXGRKCcKr9HlBEw6T9r9n/MfkHhLjo5FlhRKDa9YZRHT2ZYrnqla8Ze05fxg8rHtFd46W+9fib3HnZEFHZsoFudPpUUx79wcvnTUSIV/R2vNWPIcC2U7O3ak4HamVZowJxhVXMY/dIWaq6uSXwI4YcqM0Pe62yhx9n1VMm10URNa1F9KG6aRGPuyyKMO7JOS7z+XcGbJrdXHEMxkexUU0hfZWMcBfD6Q/SDATmdLkEhuk3CjGgAdMvRzl55JBnSefkd/oLdFCPil8jeDErg8Jb04jKCw//dHi69CtxZn7arJfEaxKWQ+WG5bBVoMIRlG1PNuZ1vtWGD6BCoxXZgmFk1qkjfDWl+/SVSQpb1N8ki5XEqud4S2BWcxZqxCRbW0sIKgnpMj5i8geMW3Z4NEbe/XNq06NwLUmwiYRJAKYYMzl7xEGbMNepegs4fBkRR0xNQbU+Mql3rLbw6nXbZbs55Z5wHnaVfe9vLURVnDGncSK1IE47XCGfFoixTtC8C4AbPm6C3NQ+nA6fQXRM2YFb0byIINi7Ej8E+s0bG2hd1aKxuNu/PtkzZw8JWhgLTxktCLELj6u9/MKyRRjjLuoKXgyQTKhEeACD87DNLQuLavZ7w1W5SUAl3HsKePqA46Lb/rUTKIUdYHgZjpSTZRrnh+wCUfkiujDp9R32Km1yeEzz3SBTkxdt+jJKUSvZSXCjbdNKUUqGeR8Os28BRbCatkZRtKAxOymWEaKhxIiRYnWYdooxFAYLpEQ0ht9RUioc6IswmFwhb45u0XjdVnswSg1Mr7qIKig0LxepqiauWNtjAIPSw1j99WbD9dYqQoVnvJ6ozpXKoPNUdLC/qPM5olCrTfzyCDvo7vvBBV4Y/hU3DuyyYFZtg/8GshGq7EPKKbVMzQD4gVokZe8LRlFcx+QfMSTwnv/3OTCatYspoUWaALzlA46TjJZ49y6w5O5f2q5m2fhXP8l/xCtJWfS/i2HXhDPoawM11ukZHE2L9IezkFwQjP1qwksM633LfPUfhNDtaHuV6uscUzwG8NlwI9kqcIJYN7Wbpst9TlawqHwgOGKujzFbpZJejt76Z5NpoiAnZhUfFqll+fgeznbMBwtVhp5NuXhM8FyDCzJCyDEg=="),
		},
	}
	for i, tc := range testcases {
		got := tc.rr.(*DNSKEY).KeyTag()
		if got != tc.tag {
			t.Errorf("test %d, expected %d, got %d", i, tc.tag, got)
		}
	}
}

func TestDNSSECVerify(t *testing.T) {
	draftKey, draftSig, draftRRs := signedRRset(t, "testdata/mldsa44-example.com", TypeMX, 59829)
	huqueKey, huqueSig, huqueRRs := signedRRset(t, "testdata/mldsa44-mldsa.huque.com", TypeDNSKEY, 23583)
	kochenKey, kochenSig, kochenRRs := signedRRset(t, "testdata/mldsa44-kochen-specker.info", TypeDNSKEY, 20767)

	testcases := []struct {
		name string
		key  *DNSKEY
		sig  *RRSIG
		rrs  []RR
	}{
		{
			"root",
			&DNSKEY{Hdr: Header{Name: ".", Class: ClassINET, TTL: 172800}, Flags: 256, Protocol: 3, Algorithm: RSASHA256,
				PublicKey: "AwEAAbauxLSFZ+KSWi2cT6TJbm3d+GIVqb2N1XnDjMsRme0b6JlGp/cvwmM5CaJ5LQ7tG1r7LuTHjYZadtbNk2nZmclq9r4KInS48ungoAZb0gJXVw8IvBTBb1YWQmiBqD285pJuORwTii7DF++nNJJk3i55HJt9SmBI7m7t8nvx7OOY/w0inxg3fLH2uY0SKO8he4FGwMc4Ubiab8N8Yhyhh+FkKKdD/+oAcuGF75PjlSXO460B4MlNLlEcjDEzIsKauRYx4YVgSaNomGhMMFblmXRzgW+1R6ywvm5mC9+omlyyizZp2GJfPwGMezuKSGDndO6CYYEc5/lsRhvBYsGjdPM=",
			},
			&RRSIG{Hdr: Header{Name: ".", Class: ClassINET, TTL: 86400}, TypeCovered: TypeSOA, Algorithm: RSASHA256, Expiration: 1760245200, Inception: 1759118400, OrigTTL: 86400, KeyTag: 46441, SignerName: ".",
				Signature: "Bi335z6iBqX1BsA6AUc29BdgoVcn2a6hfoowvC0OqNbQ1XdIz6lJ6L7iXPBYwUXTP85yII2wIBhIQJSKKuBnGmYwgXsRR7d6vxbcC63UKp6lG/xjMzae2TvuQINQhMCE2N+ufN4DCmlZspdDmWOyPkemJhncrWA+V6GGWLbCtVXPEbcgAvaNyVtFlLp788SdoDy1rpvIM2ZI9tGM2gxZG7wMLPIODtdl865H87kHaHbBkfKkCHsVCeWANoW9r8usdgt8+2HHf3w67HTERUv2NN31fgdkezaS1/7LMaMUb+5mqLygIJPVXbmK8wiUbnyuYW6Ems+R1FEnc2Dd5x5kAw==",
			},
			[]RR{&SOA{Hdr: Header{Name: ".", Class: ClassINET, TTL: 86400}, Ns: "a.root-servers.net.", Mbox: "nstld.verisign-grs.com.", Serial: 2025092900, Refresh: 1800, Retry: 900, Expire: 604800, Minttl: 86400}},
		},
		{
			"org",
			&DNSKEY{Hdr: Header{Name: "org.", Class: ClassINET, TTL: 3600}, Flags: 257, Protocol: 3, Algorithm: RSASHA256,
				PublicKey: "AwEAAexZJ/1wfyNCxNPrTZizaG7UlibGhP+AyogR6bqjptKweEgE4gD8GxRQJkt+Fn5pCoNqzmm1ZnEoKqvm93uOYtbKkYQDGH+W69J66MSKpgIyS+mT/4iaXn+lpb5o99l/sf7lHMa975O/fqN6aPUll4hUbN2T1LHv6HzQuQCtNRJA8jHGwX5q0NMmh2Z+yaG6B9cISerje9l5L+ID2ydJ6zXquYteoIUvX2xzqnXCdHPSvD+oL6R/weW+tztdFS1hok/1z3tn5NzmcaOLll9nXniCozEpLFEGPswyvtphWgCYhI8bBTqhUsIwfIwLSBQTEg2oCX7sS5CbXg44OqwhIW8=",
			},
			&RRSIG{Hdr: Header{Name: "org.", Class: ClassINET, TTL: 3600}, TypeCovered: TypeDNSKEY, Algorithm: RSASHA256, Labels: 1, Expiration: 1770478705, Inception: 1768660705, OrigTTL: 3600, KeyTag: 26974, SignerName: "org.",
				Signature: "YFVPRIIx6ZItt/17yrZmtnBQOFRD44rNDwOh3BQC+NhaM7+6w5kVVpMbVMFy6X8yuXp3+A857I6g6FYVB2p7zDRhq1hkIRyxyYKMmyQmgd0d/Km+vYU+KQjSDWFp0Cm7B+3q5bvbZKRsWho36fofO37jyWkDKYl3tPm8hSCzNuCJ7NfH+3GpcztYL/M3xeHJSJ1wwzFZUW7ioAY4cnmdzHEraXi1O/2UKb+h0lR6CdvNGSe6FLsmx1OETQ8JKneopXpm3RG07AwMSMXw+lgo1d0DZiXwscJpcqHWm9eWaI7OQX6JFl96Yjnjjh1z8gLXFYMXDxLUB3wBtYF903ukiQ==",
			},
			[]RR{
				&DNSKEY{Hdr: Header{Name: "org.", Class: ClassINET, TTL: 3600}, Flags: 256, Protocol: 3, Algorithm: RSASHA256,
					PublicKey: "AwEAAfCN1yguCELJYujNmis5cjZFEW4UcxwSitTh7m1RYMEHTjhSCCzeMZ+UrpNIxLDLwAtWCcAQPoLuQdbhbMBo17pkK2UF26k+WKLN4Ieyw8SsAFScQkf+K7wCYJ4fRJgpsvZgRYBwzpsfQjm26Jd8B7olV02AnvxUHuPo1QeGz0Dd",
				},
				&DNSKEY{Hdr: Header{Name: "org.", Class: ClassINET, TTL: 3600}, Flags: 256, Protocol: 3, Algorithm: RSASHA256,
					PublicKey: "AwEAAa0ig3MNvgtoYj85frPVdYmRm98PlHMHxQBQRRsGC4wucf5PHv9wlSQOJUfGXY45KlJpdgZX8B6s4vHZwOGW2utqp9u2JLMFaMb5bhZUlbD/qufIQu8hcLtoTkORdaK0lo+9vx8R+N2I0DcIVVWBjtSoNjb5L/ssPPicUsBobl2B",
				},
				&DNSKEY{Hdr: Header{Name: "org.", Class: ClassINET, TTL: 3600}, Flags: 257, Protocol: 3, Algorithm: RSASHA256,
					PublicKey: "AwEAAexZJ/1wfyNCxNPrTZizaG7UlibGhP+AyogR6bqjptKweEgE4gD8GxRQJkt+Fn5pCoNqzmm1ZnEoKqvm93uOYtbKkYQDGH+W69J66MSKpgIyS+mT/4iaXn+lpb5o99l/sf7lHMa975O/fqN6aPUll4hUbN2T1LHv6HzQuQCtNRJA8jHGwX5q0NMmh2Z+yaG6B9cISerje9l5L+ID2ydJ6zXquYteoIUvX2xzqnXCdHPSvD+oL6R/weW+tztdFS1hok/1z3tn5NzmcaOLll9nXniCozEpLFEGPswyvtphWgCYhI8bBTqhUsIwfIwLSBQTEg2oCX7sS5CbXg44OqwhIW8=",
				},
			},
		},
		{
			"org-from-string",
			dnstestNew("org.	30	IN	DNSKEY	257 3 8 AwEAAexZJ/1wfyNCxNPrTZizaG7UlibGhP+AyogR6bqjptKweEgE4gD8GxRQJkt+Fn5pCoNqzmm1ZnEoKqvm93uOYtbKkYQDGH+W69J66MSKpgIyS+mT/4iaXn+lpb5o99l/sf7lHMa975O/fqN6aPUll4hUbN2T1LHv6HzQuQCtNRJA8jHGwX5q0NMmh2Z+yaG6B9cISerje9l5L+ID2ydJ6zXquYteoIUvX2xzqnXCdHPSvD+oL6R/weW+tztdFS1hok/1z3tn5NzmcaOLll9nXniCozEpLFEGPswyvtphWgCYhI8bBTqhUsIwfIwLSBQTEg2oCX7sS5CbXg44OqwhIW8=").(*DNSKEY),
			dnstestNew("org.	30	IN	RRSIG	DNSKEY 8 1 3600 20260207153825 20260117143825 26974 org. YFVPRIIx6ZItt/17yrZmtnBQOFRD44rNDwOh3BQC+NhaM7+6w5kVVpMbVMFy6X8yuXp3+A857I6g6FYVB2p7zDRhq1hkIRyxyYKMmyQmgd0d/Km+vYU+KQjSDWFp0Cm7B+3q5bvbZKRsWho36fofO37jyWkDKYl3tPm8hSCzNuCJ7NfH+3GpcztYL/M3xeHJSJ1wwzFZUW7ioAY4cnmdzHEraXi1O/2UKb+h0lR6CdvNGSe6FLsmx1OETQ8JKneopXpm3RG07AwMSMXw+lgo1d0DZiXwscJpcqHWm9eWaI7OQX6JFl96Yjnjjh1z8gLXFYMXDxLUB3wBtYF903ukiQ==").(*RRSIG),
			[]RR{
				dnstestNew("org.	30	IN	DNSKEY	256 3 8 AwEAAfCN1yguCELJYujNmis5cjZFEW4UcxwSitTh7m1RYMEHTjhSCCzeMZ+UrpNIxLDLwAtWCcAQPoLuQdbhbMBo17pkK2UF26k+WKLN4Ieyw8SsAFScQkf+K7wCYJ4fRJgpsvZgRYBwzpsfQjm26Jd8B7olV02AnvxUHuPo1QeGz0Dd"),
				dnstestNew("org.	30	IN	DNSKEY	256 3 8 AwEAAa0ig3MNvgtoYj85frPVdYmRm98PlHMHxQBQRRsGC4wucf5PHv9wlSQOJUfGXY45KlJpdgZX8B6s4vHZwOGW2utqp9u2JLMFaMb5bhZUlbD/qufIQu8hcLtoTkORdaK0lo+9vx8R+N2I0DcIVVWBjtSoNjb5L/ssPPicUsBobl2B"),
				dnstestNew("org.	30	IN	DNSKEY	257 3 8 AwEAAexZJ/1wfyNCxNPrTZizaG7UlibGhP+AyogR6bqjptKweEgE4gD8GxRQJkt+Fn5pCoNqzmm1ZnEoKqvm93uOYtbKkYQDGH+W69J66MSKpgIyS+mT/4iaXn+lpb5o99l/sf7lHMa975O/fqN6aPUll4hUbN2T1LHv6HzQuQCtNRJA8jHGwX5q0NMmh2Z+yaG6B9cISerje9l5L+ID2ydJ6zXquYteoIUvX2xzqnXCdHPSvD+oL6R/weW+tztdFS1hok/1z3tn5NzmcaOLll9nXniCozEpLFEGPswyvtphWgCYhI8bBTqhUsIwfIwLSBQTEg2oCX7sS5CbXg44OqwhIW8="),
			},
		},
		{"mldsa44-example.com", draftKey, draftSig, draftRRs},
		{"mldsa44-mldsa.huque.com", huqueKey, huqueSig, huqueRRs},
		{"mldsa44-kochen-specker.info", kochenKey, kochenSig, kochenRRs},
	}

	options := &SignOption{}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sig.Verify(tc.key, tc.rrs, options)
			if err != nil {
				t.Fatalf("failure to verify: %s", err)
			}
		})
	}
}

func signedRRset(t *testing.T, file string, typ, keytag uint16) (*DNSKEY, *RRSIG, []RR) {
	t.Helper()
	var key *DNSKEY
	var sig *RRSIG
	rrs := []RR{}
	for _, rr := range readZone(t, file) {
		switch x := rr.(type) {
		case *DNSKEY:
			if x.KeyTag() == keytag {
				key = x
			}
		case *RRSIG:
			if x.TypeCovered == typ && x.KeyTag == keytag {
				sig = x
			}
		}
		if RRToType(rr) == typ {
			rrs = append(rrs, rr)
		}
	}
	if key == nil || sig == nil {
		t.Fatalf("no DNSKEY or RRSIG with keytag %d in %s", keytag, file)
	}
	return key, sig, rrs
}

// mldsa44PrivateKey is the private key from section 6 of draft-westerbaan-dnssec-mldsa.
const mldsa44PrivateKey = `Private-key-format: v1.3
Algorithm: 18 (MLDSA44)
PrivateKey: AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8=
`

func TestDNSSECMLDSA44PrivateKey(t *testing.T) {
	var key *DNSKEY
	for _, rr := range readZone(t, "testdata/mldsa44-example.com") {
		if k, ok := rr.(*DNSKEY); ok {
			key = k
		}
	}

	priv, err := key.NewPrivate(mldsa44PrivateKey)
	if err != nil {
		t.Fatalf("failure to read the private key: %s", err)
	}

	// The key pair is derived from the seed, the public part must match the DNSKEY.
	pub := new(DNSKEY)
	pub.setPublicKeyMLDSA44(priv.(*mldsa.PrivateKey).PublicKey())
	if pub.PublicKey != key.PublicKey {
		t.Error("public key derived from the seed does not match the DNSKEY")
	}

	if got := key.PrivateKeyString(priv); got != mldsa44PrivateKey {
		t.Errorf("expected private key %q, got %q", mldsa44PrivateKey, got)
	}
}
