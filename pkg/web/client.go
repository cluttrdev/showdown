package web

import (
	"golang.org/x/net/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager

	egress chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return *Client{
		connection: conn,
		manager:    manager,
		egress: make(chan Message),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}

	for {
		var msg Message
		if err := websocket.JSON.Receive(c.connection, &msg); err != nil {
			log.Printf("readMessages: %v", err)
			break
		}

		log.Printf("Received message: %s", string(msg.data))
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}

	for {
		select {
		case msg, ok := <-c.egress:
			// ok will be false in case the egress channel is closed
			if !ok {
				// return to close the goroutine
				return
			}

			if err := websocket.JSON.Send(c.connection, msg); err != nil {
				log.Printf("writeMessages: %v", err)
			}
		}
	}
}
