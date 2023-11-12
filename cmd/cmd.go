package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/cluttrdev/showdown/internal/command"
	"github.com/cluttrdev/showdown/internal/version"
)

var rootCmd *command.Command

func Configure(v version.Info) {
	var (
		stopCmd    = NewStopCmd()
		versionCmd = NewVersionCmd(v)
		termCmd    = NewTerminalCmd()
	)

	rootCmd = NewRootCmd()

	rootCmd.Subcommands = []*command.Command{
		stopCmd,
		termCmd,
		versionCmd,
	}
}

func Execute() error {
	if rootCmd == nil {
		return errors.New("not configured")
	}

	if err := rootCmd.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		} else {
			return fmt.Errorf("error parsing arguments: %w", err)
		}
	}

	ctx := context.Background()
	return rootCmd.Run(ctx)
}
