package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gomarkdown/markdown"

	watch "github.com/cluttrdev/showdown/internal"
)

type application struct {
	FilePath string
	Content  []byte
}

func main() {
	file := flag.String("file", "", "The file to preview")
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	if *file == "" {
		log.Fatal("No file given!")
	}

	app := &application{
		FilePath: *file,
	}

	app.render()

	handler := func(e fsnotify.Event) {
		if e.Has(fsnotify.Write) {
			app.render()
		}
	}

	w, err := watch.WatchFile(*file, handler)
	defer w.Close()

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.root)

	return mux
}

func (app *application) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Write(app.Content)
}

func (app *application) render() {
	in, err := os.ReadFile(app.FilePath)
	if err != nil {
		return
	}

	app.Content = markdown.ToHTML(in, nil, nil)
}
