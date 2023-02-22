package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
)

type application struct {
	FilePath string
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

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
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

	app.render(w)
}

func (app *application) render(w http.ResponseWriter) {
	in, err := os.ReadFile(app.FilePath)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	out := markdown.ToHTML(in, nil, nil)

	w.Write(out)
}
