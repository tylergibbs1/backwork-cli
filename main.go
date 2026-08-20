package main

import (
	"fmt"
	"os"

	"github.com/tylergibbs1/backwork-cli/cmd"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	cmd.SetBuildInfo(Version, BuildTime)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
