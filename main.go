package main

import (
	"fmt"
	"os"

	"github.com/cluttrdev/showdown/cmd"

	versionpkg "github.com/cluttrdev/showdown/internal/version"
)

var (
	version = "(devel)"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Configure(versionpkg.Info{
		Version: version,
		Commit:  commit,
		Date:    date,
	})

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
