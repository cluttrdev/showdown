package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
	"golang.org/x/net/websocket"

	watch "github.com/cluttrdev/showdown/internal"
)

type Application struct {
	file    string
	sockets map[*websocket.Conn]bool
}

func NewApplication(file string) *Application {
	return &Application{
		file:    file,
		sockets: make(map[*websocket.Conn]bool),
	}
}

func (app *Application) run() error {
	w, err := watch.WatchFile(app.file, func() {
		content, err := app.render()
		if err != nil {
			log.Fatal(err)
			return
		}

		for ws := range app.sockets {
			app.sendContent(ws, string(content))
		}
	})
	if err != nil {
		return err
	}
	defer w.Close()

	addr := fmt.Sprintf("127.0.0.1:%s", "1337")
	return app.serve(addr)
}

func (app *Application) serve(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}

	return srv.ListenAndServe()
}

func (app *Application) render() ([]byte, error) {
	in, err := os.ReadFile(app.file)
	if err != nil {
		return nil, err
	}

	out := markdown.ToHTML(in, nil, nil)

	return out, nil
}
