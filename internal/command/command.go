package command

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

type Command struct {
	Name       string
	ShortHelp  string
	ShortUsage string
	LongHelp   string

	Flags *flag.FlagSet
	Exec  func(ctx context.Context, args []string) error
}

func (cmd *Command) Parse(args []string) error {
	if cmd.Name == "" {
		return errors.New("name is required")
	}
	if cmd.Flags == nil {
		cmd.Flags = flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
	}

	cmd.Flags.Usage = func() {
		fmt.Fprintln(cmd.Flags.Output(), usageFunc(cmd))
	}

	if err := cmd.Flags.Parse(args); err != nil {
		return err
	}

	return nil
}

func (cmd *Command) Run(ctx context.Context) (err error) {
	if !cmd.Flags.Parsed() {
		return errors.New("not parsed")
	}

	if cmd.Exec == nil {
		return fmt.Errorf("%s: %w", cmd.Name, errors.New("no exec function"))
	}

	defer func() {
		if errors.Is(err, flag.ErrHelp) {
			cmd.Flags.Usage()
			err = nil
		}
	}()

	if err := cmd.Exec(ctx, cmd.Flags.Args()); err != nil {
		return err
	}
	return nil
}
