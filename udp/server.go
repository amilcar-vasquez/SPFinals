package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func StartUDPServer() {
	addr := net.UDPAddr{
		Port: 9001,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error starting UDP server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("UDP server started on port", addr.Port)

	go cleanupInactiveClients()

	buffer := make([]byte, 1024)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		msg := strings.TrimSpace(string(buffer[:n]))
		incrementMessageCount(n)
		clientKey := clientAddr.String()

		mu.Lock()
		client, exists := clients[clientKey]
		if !exists {
			client = Client{Addr: clientAddr}
		}
		client.LastSeen = time.Now()

		if strings.HasPrefix(msg, "/name ") {
			name := strings.TrimSpace(msg[len("/name "):])
			if name != "" {
				client.Name = name
				clients[clientKey] = client
				conn.WriteToUDP([]byte("✅ Name set to: "+name), clientAddr)
				mu.Unlock()
				continue
			}
		}

		clients[clientKey] = client
		mu.Unlock()

		sender := client.Name
		if sender == "" {
			sender = clientKey
		}

		// If it's a latency message, don't wrap it — broadcast as-is
if strings.HasPrefix(msg, "[LATENCY]|") {
	go broadcast(conn, msg, clientAddr)
} else {
	go broadcast(conn, fmt.Sprintf("[%s]: %s", sender, msg), clientAddr)
}

	}
}
