package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

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
				if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
					// websocket connection closed by either client or server
					delete(s.sockets, ws)
					break
				}
				log.Printf("Receive error: %v", err)
			}
		}
	}()

	if err := sendMessage(ws, MessageTypeTitle, s.Title); err != nil {
		log.Printf("error sending title: %v\n", err)
	}

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
	MessageTypeTitle   = MessageType("title")
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
