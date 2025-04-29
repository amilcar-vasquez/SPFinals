// file: clientHandlers.go

package main

import (
	"net"
	"sync"
)

type Client struct {
	conn     net.Conn
	nickname string
}

type ClientManager struct {
	clients map[net.Conn]*Client
	mu      sync.Mutex
}

var manager = ClientManager{
	clients: make(map[net.Conn]*Client),
}

// Add a client to the manager
func (m *ClientManager) Add(conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[conn] = &Client{
		conn:     conn,
		nickname: "", // Default empty nickname
	}
}

// Remove a client from the manager
func (m *ClientManager) Remove(conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, conn)
}

// Broadcast message to all clients except sender
func (m *ClientManager) Broadcast(sender net.Conn, message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for conn, client := range m.clients {
		if conn != sender {
			client.conn.Write([]byte(message + "\n"))
		}
	}
}

// SetNickname updates nickname and returns the old one
func (m *ClientManager) SetNickname(conn net.Conn, newNick string) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.clients[conn]
	if !exists {
		return "Unknown"
	}
	oldNick := client.nickname
	client.nickname = newNick
	return oldNick
}

// GetClient returns a pointer to the client struct
func (m *ClientManager) GetClient(conn net.Conn) *Client {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.clients[conn]
}
