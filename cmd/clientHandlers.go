// file: cmd/clientHandlers.go

package main

import (
	"github.com/fatih/color"
	"net"
	"sync"
)

type Client struct {
	conn      net.Conn
	nickname  string
	colorFunc func(format string, a ...interface{}) string
}

type ClientManager struct {
	clients map[net.Conn]*Client
	mu      sync.Mutex
}

var manager = ClientManager{
	clients: make(map[net.Conn]*Client),
}

var colorFuncs = []func(string, ...interface{}) string{
	color.HiBlueString,
	color.HiGreenString,
	color.HiRedString,
	color.HiMagentaString,
	color.HiCyanString,
	color.HiWhiteString,
}

// Add a new client to the manager
func (m *ClientManager) Add(conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Pick a color using round-robin or hash on IP
	colorFunc := colorFuncs[len(m.clients)%len(colorFuncs)]

	m.clients[conn] = &Client{
		conn:      conn,
		nickname:  "",
		colorFunc: colorFunc,
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
	// Locking the manager to ensure thread safety
	m.mu.Lock()
	// Create a copy of the clients slice to avoid concurrent map access
	clientsCopy := make([]*Client, 0, len(m.clients))
	senderClient := m.clients[sender]
	for _, client := range m.clients {
		if client.conn != sender {
			clientsCopy = append(clientsCopy, client)
		}
	}
	m.mu.Unlock()

	colored := senderClient.colorFunc(message)
	for _, client := range clientsCopy {
		n, err := client.conn.Write([]byte(colored + "\n"))
		if err != nil {
			metrics.IncrementFailedWrite(client.conn)
		} else {
			metrics.AddBytes(client.conn, n)
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
