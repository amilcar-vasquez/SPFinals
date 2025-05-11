// file: metrics.go
package main

import (
	"fmt"
	"github.com/fatih/color"
	"net"
	"sync"
	"time"
)

type ClientMetrics struct {
	MessageCount    int
	ConnectionTime  time.Time
	TotalBytes      int
	LatencySamples  []time.Duration
	FailedWrites    int
	LastPingSent    time.Time
	LastPingLatency time.Duration
}

type Metrics struct {
	mu            sync.Mutex
	TotalMessages int
	ActiveClients int
	ClientsData   map[string]*ClientMetrics
}

var metrics = Metrics{
	ClientsData: make(map[string]*ClientMetrics),
}

// Call on new connection
func (m *Metrics) RegisterClient(conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	addr := conn.RemoteAddr().String()
	m.ClientsData[addr] = &ClientMetrics{
		ConnectionTime: time.Now(),
	}
	m.ActiveClients++
}

// Call on disconnection
func (m *Metrics) UnregisterClient(conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	addr := conn.RemoteAddr().String()
	delete(m.ClientsData, addr)
	m.ActiveClients--
}

// Call on each message
func (m *Metrics) IncrementMessage(conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	addr := conn.RemoteAddr().String()
	m.TotalMessages++
	if client, ok := m.ClientsData[addr]; ok {
		client.MessageCount++
	}
}

// Return current snapshot
func (m *Metrics) GetSnapshot() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := color.New(color.FgCyan).Sprint("=== Server Metrics ===\n")
	result += color.New(color.FgGreen).Sprintf("Total Messages: %d\n", m.TotalMessages)
	result += color.New(color.FgYellow).Sprintf("Active Clients: %d\n", m.ActiveClients)
	result += color.New(color.FgMagenta).Sprint("Client Breakdown:\n")

	totalLatency := time.Duration(0)
	latencySamples := 0

	for addr, data := range m.ClientsData {
		duration := time.Since(data.ConnectionTime).Seconds()
		msgs := data.MessageCount
		throughput := 0.0
		if duration > 0 {
			throughput = float64(data.TotalBytes) / duration
		}

		packetLoss := data.FailedWrites

		// Latency aggregation
		for _, sample := range data.LatencySamples {
			totalLatency += sample
			latencySamples++
		}

		result += fmt.Sprintf(" - %s:\n", addr)
		result += fmt.Sprintf("    Messages: %d\n", msgs)
		result += fmt.Sprintf("    Connected: %.0f sec\n", duration)
		result += fmt.Sprintf("    Throughput: %.2f bytes/sec\n", throughput)
		result += fmt.Sprintf("    Packet Loss: %d failed writes\n", packetLoss)
		// You could even show most recent latency here if you want:
		result += fmt.Sprintf("    Last Latency: %v\n", data.LastPingLatency)
	}

	avgLatency := time.Duration(0)
	if latencySamples > 0 {
		avgLatency = totalLatency / time.Duration(latencySamples)
	}

	result += fmt.Sprintf("Average Latency: %v\n", avgLatency)

	return result
}

func (m *Metrics) AddBytes(conn net.Conn, n int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	addr := conn.RemoteAddr().String()
	if client, ok := m.ClientsData[addr]; ok {
		client.TotalBytes += n
	}
}

func (m *Metrics) IncrementFailedWrite(conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	addr := conn.RemoteAddr().String()
	if client, ok := m.ClientsData[addr]; ok {
		client.FailedWrites++
	}
}

func (m *Metrics) RecordPingSent(conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	addr := conn.RemoteAddr().String()
	if client, ok := m.ClientsData[addr]; ok {
		client.LastPingSent = time.Now()
	}
}

func (m *Metrics) RecordPingLatency(conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	addr := conn.RemoteAddr().String()
	if client, ok := m.ClientsData[addr]; ok && !client.LastPingSent.IsZero() {
		client.LastPingLatency = time.Since(client.LastPingSent)
		client.LatencySamples = append(client.LatencySamples, client.LastPingLatency)
	}
}
