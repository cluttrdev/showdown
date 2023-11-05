package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:generate cp -r ../../web/static/ ./assets
//go:embed assets/index.html
//go:embed assets/js/client.js
var assets embed.FS

func fileServer() (http.Handler, error) {
    fsys, err := fs.Sub(assets, "assets")
    if err != nil {
        return nil, fmt.Errorf("error creating filesystem subtree: %w", err)
    }

    return http.FileServer(http.FS(fsys)), nil
}
