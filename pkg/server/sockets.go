package server

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

func (s *Server) handleWebSocket(ws *websocket.Conn) {
    defer ws.Close()
    
    ctx, cancel := context.WithCancel(context.Background())
    s.sockets[ws] = cancel

    go func() {
        defer cancel()

        for {
            var msg []byte
            if err := websocket.Message.Receive(ws, &msg); err != nil {
                log.Printf("Receive error: %+v", err)
                delete(s.sockets, ws)
                break
            }
        }
    }()

    // trigger sending content
    s.Update()

    <-ctx.Done()
}

func (s *Server) closeWebSockets() {
    for _, cancel := range s.sockets {
        cancel()
    }
}

type MessageType string

const (
    MessageTypeTitle = MessageType("title")
    MessageTypeContent = MessageType("content")
)

type Message struct {
    Type string `json:"type"`
    Data string `json:"data"`
}

func sendMessage(ws *websocket.Conn, t MessageType, d string) error {
    msg := Message{
        Type: string(t),
        Data: d,
    }

    if err := websocket.JSON.Send(ws, msg); err != nil {
        return fmt.Errorf("error sending message: %w", err)
    }

    return nil
}

func sendTitle(ws *websocket.Conn, title string) error {
    return sendMessage(ws, MessageTypeTitle, title)
}

func sendContent(ws *websocket.Conn, content string) error {
    return sendMessage(ws, MessageTypeContent, content)
}
