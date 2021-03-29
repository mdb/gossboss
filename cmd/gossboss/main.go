package main

import (
	"github.com/mdb/gossboss/cmd"
)

// version's value is passed in at build time via -ldflags "-X main.version=$VERSION".
// By default, goreleaser passes:
// `-s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}} -X main.builtBy=goreleaser`.
var version string

func main() {
	cmd.Execute(version)
}
