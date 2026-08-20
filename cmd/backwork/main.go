package main

import (
	"fmt"
	"os"

	root "github.com/tylergibbs1/backwork-cli/cmd"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	root.SetBuildInfo(Version, BuildTime)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
