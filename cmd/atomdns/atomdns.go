package main

import (
	"github.com/lanrat/dns/cmd/atomdns/atom"
)

//go:generate go run man_generate.go
//go:generate go run release_generate.go

const version = "075"

func main() { atom.Run(version) }
