package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

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
