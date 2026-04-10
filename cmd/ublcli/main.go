package main

import (
	"ublcli/internal/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.Run(cli.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	})
}
