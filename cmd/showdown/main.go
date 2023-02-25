package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"

	"github.com/gomarkdown/markdown"
	"github.com/pkg/errors"

	watch "github.com/cluttrdev/showdown/internal"
)

type Application struct {
	FilePath string

	sockets map[*websocket.Conn]bool
}

func main() {
	file := flag.String("file", "", "The file to preview")
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	if *file == "" {
		log.Fatal("No file given!")
		return
	}

	app := &Application{
		FilePath: *file,
		sockets:  make(map[*websocket.Conn]bool),
	}

	w, err := watch.WatchFile(app.FilePath, func() {
		for ws := range app.sockets {
			app.sendContent(ws)
		}
	})
	defer w.Close()

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	log.Fatal(errors.Errorf("Stopped listening: %v", err))
}

func (app *Application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.root)
	mux.Handle("/ws", websocket.Handler(app.socket))

	return mux
}

func (app *Application) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, "./web/index.html")
}

func (app *Application) socket(ws *websocket.Conn) {
	if err := app.sendTitle(ws); err != nil {
		log.Fatal(err)
		return
	}

	if err := app.sendContent(ws); err != nil {
		log.Fatal(err)
		return
	}

	defer ws.Close()
	app.sockets[ws] = true

	for {
		var msg []byte
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			log.Printf("Receive: %v", err)
			delete(app.sockets, ws)
			break
		}
	}
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func (app *Application) sendTitle(ws *websocket.Conn) error {
	msg := Message{
		Type: "title",
		Data: app.FilePath,
	}

	if err := websocket.JSON.Send(ws, msg); err != nil {
		return errors.Errorf("sendTitle: %v", err)
	}

	return nil
}

func (app *Application) sendContent(ws *websocket.Conn) error {
	content, err := app.render()
	if err != nil {
		return errors.Errorf("sendContent / render: %v", err)
	}

	msg := Message{
		Type: "content",
		Data: string(content),
	}

	if err = websocket.JSON.Send(ws, msg); err != nil {
		return errors.Errorf("sendContent / send: %v", err)
	}
	return nil
}

func (app *Application) render() ([]byte, error) {
	in, err := os.ReadFile(app.FilePath)
	if err != nil {
		return nil, err
	}

	out := markdown.ToHTML(in, nil, nil)

	return out, nil
}
