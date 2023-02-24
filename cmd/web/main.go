package main

import (
	"bytes"
	"flag"
	"html/template"
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
	Template *template.Template
}

func main() {
	file := flag.String("file", "", "The file to preview")
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	if *file == "" {
		log.Fatal("No file given!")
		return
	}

	t, err := template.ParseFiles("./web/templates/base.html.tpl")
	if err != nil {
		log.Fatal(err)
		return
	}

	app := &application{
		FilePath: *file,
		Template: t,
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

	var data struct {
		Content template.HTML
	}
	data.Content = template.HTML(app.Content)

	buf := new(bytes.Buffer)

	err := app.Template.ExecuteTemplate(buf, "base", data)
	if err != nil {
		log.Fatal(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	buf.WriteTo(w)
}

func (app *application) render() {
	in, err := os.ReadFile(app.FilePath)
	if err != nil {
		return
	}

	app.Content = markdown.ToHTML(in, nil, nil)
}
