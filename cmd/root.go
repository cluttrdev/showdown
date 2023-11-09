package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/cluttrdev/showdown/pkg/content"
	"github.com/cluttrdev/showdown/pkg/server"

	"github.com/cluttrdev/showdown/internal/command"
	"github.com/cluttrdev/showdown/internal/watch"
)

type RootCmdConfig struct {
	host string
	port string
}

func NewRootCmd() *command.Command {
	cfg := RootCmdConfig{}

	fs := flag.NewFlagSet("showdown", flag.ContinueOnError)

	cfg.RegisterFlags(fs)

	return &command.Command{
		Name:       "showdown",
		ShortHelp:  "showdown - A live markdown previewer",
		ShortUsage: "showdown [flags] FILE",
		LongHelp:   "",
		Flags:      fs,
		Exec:       cfg.Exec,
	}
}

func (c *RootCmdConfig) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.host, "host", "127.0.0.1", "The address the server listens to.")
	fs.StringVar(&c.port, "port", "1337", "The port the server listens on.")
}

func (c *RootCmdConfig) Exec(ctx context.Context, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: not enough arguments")
		return flag.ErrHelp
	} else if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Error: too many arguments")
		return flag.ErrHelp
	}

	return c.run(ctx, args[0])
}

func (c *RootCmdConfig) run(ctx context.Context, file string) error {
	r := &content.MarkdownRenderer{
		File: file,
	}

	srv := server.Server{
		Title:    r.File,
		Renderer: r,
	}

	g, ctx := errgroup.WithContext(ctx)

	// start server
	g.Go(func() error {
		addr := fmt.Sprintf("%s:%s", c.host, c.port)
		return srv.Serve(ctx, addr)
	})

	// watch file for changes
	g.Go(func() error {
		w, err := watch.WatchFile(file, srv.Update)
		if err != nil {
			return err
		}
		defer w.Close()

		<-ctx.Done()
		return ctx.Err()
	})

	// wait for signal to exit
	g.Go(func() error {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-signalChan:
			signal.Stop(signalChan)
			return errors.New("Got SIGINT/SIGTERM, exiting")
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	if err := g.Wait(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
