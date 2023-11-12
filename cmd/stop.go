package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/cluttrdev/showdown/internal/command"
)

type StopCmdConfig struct {
	RootCmdConfig
}

func NewStopCmd() *command.Command {
	cfg := StopCmdConfig{}

	fs := flag.NewFlagSet("stop", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "stop",
		ShortHelp:  "Stop the preview server",
		ShortUsage: "showdown stop [--port PORT]",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

func (c *StopCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootCmdConfig.RegisterFlags(fs)
}

func (c *StopCmdConfig) Exec(ctx context.Context, args []string) error {
	url := fmt.Sprintf("http://%s:%s/shutdown", c.host, c.port)

	res, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return fmt.Errorf("error requesting server stop: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("error reading response: %v\n", err)
		}
		return fmt.Errorf("error stopping server: %v", string(body))
	}
	return nil
}
