package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"text/template"
)

//go:generate mkdir -p assets
//go:generate cp -r ../../web/static/ ./assets/static
//go:generate cp -r ../../web/templates/ ./assets/templates

var (
	//go:embed assets/static
	staticFS  embed.FS
	staticDir string = "assets/static"

	//go:embed assets/templates
	templatesFS  embed.FS
	templatesDir string = "assets/templates"
)

func staticFileServer() (http.Handler, error) {
	fsys, err := fs.Sub(staticFS, staticDir)
	if err != nil {
		return nil, fmt.Errorf("error creating filesystem subtree: %w", err)
	}

	return http.FileServer(http.FS(fsys)), nil
}

type Style struct {
	Name    string
	Variant string
}

func rootHandler(s string) http.Handler {
	unpack := func(vals []string, vars ...*string) {
		for i, val := range vals {
			*vars[i] = val
		}
	}

	var style Style
	unpack(strings.Split(s, ":"), &style.Name, &style.Variant)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		templateFiles := []string{
			templatesDir + "/" + "index.html",
			templatesDir + "/" + "styles" + "/" + style.Name + ".html",
		}

		ts, err := template.ParseFS(templatesFS, templateFiles...)
		if err != nil {
			log.Println(err.Error())
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}

		data := struct {
			Style Style
		}{
			Style: style,
		}
		if err := ts.Execute(w, data); err != nil {
			log.Println(err.Error())
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
	})
}
