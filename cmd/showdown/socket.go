package main

import (
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func (app *Application) sendTitle(ws *websocket.Conn, title string) error {
	msg := Message{
		Type: "title",
		Data: title,
	}

	if err := websocket.JSON.Send(ws, msg); err != nil {
		return errors.Errorf("sendTitle: %v", err)
	}

	return nil
}

func (app *Application) sendContent(ws *websocket.Conn, content string) error {
	msg := Message{
		Type: "content",
		Data: content,
	}

	if err := websocket.JSON.Send(ws, msg); err != nil {
		return errors.Errorf("sendContent / send: %v", err)
	}
	return nil
}
