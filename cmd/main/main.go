package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"

	"github.com/fsnotify/fsnotify"
	"github.com/gomarkdown/markdown"
	"github.com/pkg/errors"

	watch "github.com/cluttrdev/showdown/internal"
)

type application struct {
	FilePath string
	manager  Manager
}

func main() {
	file := flag.String("file", "", "The file to preview")
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	if *file == "" {
		log.Fatal("No file given!")
		return
	}

	manager := NewManager()

	app := &application{
		FilePath: *file,
		manager:  manager,
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./web/"))

	mux.Handle("/assets/", http.StripPrefix("/assets", fileServer))

	mux.HandleFunc("/", app.root)
	mux.Handle("/ws", websocket.Handler(app.manager.serveWS))

	return mux
}

func (app *application) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, "./web/index.html")
}

func (app *application) socket(ws *websocket.Conn) {
	if err := app.sendTitle(ws); err != nil {
		log.Fatal(err)
		return
	}

	if err := app.sendContent(ws); err != nil {
		log.Fatal(err)
		return
	}

	handler := func(e fsnotify.Event) {
		if e.Has(fsnotify.Write) {
			if err := app.sendContent(ws); err != nil {
				log.Fatal(err)
			}
		}
	}

	watch.WatchFile(app.FilePath, handler)
	log.Println("STOPPED WATCHING")
}

func (app *application) sendTitle(ws *websocket.Conn) error {
	msg := Message{
		Type: "title",
		Data: app.FilePath,
	}

	if err := websocket.JSON.Send(ws, msg); err != nil {
		return errors.Errorf("sendTitle: %v", err)
	}

	return nil
}

func (app *application) sendContent(ws *websocket.Conn) error {
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

func (app *application) render() ([]byte, error) {
	in, err := os.ReadFile(app.FilePath)
	if err != nil {
		return nil, err
	}

	out := markdown.ToHTML(in, nil, nil)

	return out, nil
}
