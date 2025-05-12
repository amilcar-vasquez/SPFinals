package main

import (
	"fmt"
	"math/rand"
	"net"
)

func broadcast(conn *net.UDPConn, message string, sender *net.UDPAddr) {
	mu.Lock()
	defer mu.Unlock()

	for addrStr, client := range clients {
		if addrStr == sender.String() {
			continue
		}

		// Simulate packet loss
		if rand.Float64() < 0.10 {
			fmt.Println("⚠️ Simulating packet loss to", addrStr)
			incrementSimulatedDrops()
			continue
		}

		_, err := conn.WriteToUDP([]byte(message), client.Addr)
		incrementBroadcastCount()

		if err != nil {
			fmt.Printf("❌ Failed to send to %s: %v. Removing client.\n", addrStr, err)
			delete(clients, addrStr)
		}
	}
}

