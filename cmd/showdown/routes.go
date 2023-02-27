package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

//go:generate cp -r ../../web/static/ ./assets
//go:embed assets/index.html
//go:embed assets/js/client.js
var assets embed.FS

func (app *Application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	sub, err := fs.Sub(assets, "assets")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(sub))
	mux.Handle("/", fileServer)
	mux.Handle("/ws", websocket.Handler(app.socket))

	return mux
}

func (app *Application) socket(ws *websocket.Conn) {
	if err := app.sendTitle(ws, app.file); err != nil {
		log.Fatal(err)
		return
	}

	content, err := app.render()
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := app.sendContent(ws, string(content)); err != nil {
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
