package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
)

func Execute() error {
	var (
		rootCmd = NewRootCmd()
		stopCmd = NewStopCmd()
	)

	rootCmd.Subcommands = append(rootCmd.Subcommands, stopCmd)

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
