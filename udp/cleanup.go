package main

import (
	"fmt"
	"time"
)

func cleanupInactiveClients() {
	for {
		time.Sleep(10 * time.Second)

		fmt.Println("🧹 Cleaning up inactive clients...") // <--- Add this log line

		mu.Lock()
		for addrStr, client := range clients {
			if time.Since(client.LastSeen) > 1*time.Minute {
				fmt.Println("Removing inactive client:", addrStr)
				delete(clients, addrStr)
			}
		}
		mu.Unlock()
	}
}
