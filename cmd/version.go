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

type VersionCmdConfig struct {
	version version.Info

	showCommit bool
	showDate   bool
}

func NewVersionCmd(v version.Info) *command.Command {
	cfg := VersionCmdConfig{
		version: v,
	}

	fs := flag.NewFlagSet("version", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "version",
		ShortHelp:  "Show version information",
		ShortUsage: "showdown version [option]...",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

func (c *VersionCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	fs.BoolVar(&c.showCommit, "commit", false, "Show git commit SHA.")
	fs.BoolVar(&c.showDate, "date", false, "Show git commit date.")
}

func (c *VersionCmdConfig) Exec(ctx context.Context, args []string) error {
	if len(args) > 0 {
		fmt.Fprintln(os.Stderr, "Error: too many arguments")
		return flag.ErrHelp
	}

	if c.showCommit && c.showDate {
		return errors.New("version: flags --commit and --date are mutual exclusive")
	} else if c.showCommit {
		fmt.Println(c.version.Commit)
	} else if c.showDate {
		fmt.Println(c.version.Date)
	} else {
		fmt.Println(c.version.Version)
	}

	return nil
}
