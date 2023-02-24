package web

import (
	"log"

	"golang.org/x/net/websocket"
)

type Manager struct {
	clients ClientList
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

func (m *Manager) serveWS(conn *websocket.Conn) {
	log.Println("New connection")

	c := NewClient(conn, m)
	m.addClient(c)

	go c.readMessages()
	go c.writeMessages()
}

func (m *Manager) addClient(c *Client) {
	m.clients[c] = true
}

func (m *Manager) removeClient(c *Client) {
	if _, ok := m.clients[c]; ok {
		c.connection.Close()
		delete(m.clients, c)
	}
}
