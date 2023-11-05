package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cluttrdev/showdown/pkg/content"
	"github.com/cluttrdev/showdown/pkg/server"

	watch "github.com/cluttrdev/showdown/internal"
)

type Application struct {
    file string
}

func NewApplication(file string) *Application {
	return &Application{
        file: file,
	}
}

func (app *Application) run(port uint16) error {
    r := &content.MarkdownRenderer{
        File: app.file,
    }

    srv := server.Server{
        Title: r.File,
        Renderer: r,
    }

	w, err := watch.WatchFile(r.File, srv.Update)
	if err != nil {
		return err
	}
	defer w.Close()

	addr := fmt.Sprintf("127.0.0.1:%d", port)
    if err := srv.Serve(context.Background(), addr); err != http.ErrServerClosed {
        return err
    }
    return nil
}

func StopServer(port uint16) error {
	url := fmt.Sprintf("http://127.0.0.1:%d/shutdown", port)
	_, err := http.Get(url)
	return err
}
