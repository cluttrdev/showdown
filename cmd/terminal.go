package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/cluttrdev/showdown/pkg/content"

	"github.com/cluttrdev/showdown/internal/command"
)

type TerminalCmdConfig struct {
	lineWidth int
	leftPad   int

	output string
}

func NewTerminalCmd() *command.Command {
	cfg := TerminalCmdConfig{}

	fs := flag.NewFlagSet("terminal", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "terminal",
		ShortHelp:  "Render markdown for terminal",
		ShortUsage: "showdown terminal [option]... FILE",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

func (c *TerminalCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	fs.IntVar(&c.lineWidth, "line-width", 80, "Maximum line width allowed.")
	fs.IntVar(&c.leftPad, "left-pad", 2, "Constant left padding.")

	fs.StringVar(&c.output, "output", "-", "Write output to file instead of stdout.")
}

func (c *TerminalCmdConfig) Exec(ctx context.Context, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: not enough arguments")
		return flag.ErrHelp
	} else if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Error: too many arguments")
		return flag.ErrHelp
	}

	file := args[0]

	r := &content.TerminalRenderer{
		File:    file,
		Width:   c.lineWidth,
		LeftPad: c.leftPad,
	}
	result, err := r.Render()
	if err != nil {
		return err
	}

	if c.output != "-" {
		if err := os.WriteFile(c.output, result, 0666); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	} else {
		_, err := fmt.Println(string(result))
		if err != nil {
			return fmt.Errorf("error printing rendered content: %w", err)
		}
	}

	return nil
}
