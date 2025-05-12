package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	serverAddr := "127.0.0.1:9001"
	raddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("Failed to resolve address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to UDP server at", serverAddr)
	fmt.Println("Type your messages below. Use '/name YourName' to set your name.")



	// Listen for replies
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println("Error reading from server:", err)
				return
			}
			received := string(buffer[:n])

			// Check if message has timestamp tag
			if strings.HasPrefix(received, "[LATENCY]") {
				parts := strings.SplitN(received, "|", 3)
				if len(parts) == 3 {
					sentTimeStr := parts[1]
					msg := parts[2]
					sentTime, err := time.Parse(time.RFC3339Nano, sentTimeStr)
					if err == nil {
						latency := time.Since(sentTime)
						fmt.Printf("📨 %s (RTT: %s)\n", msg, latency.Truncate(time.Microsecond))
						continue
					}
				}
			}

			fmt.Println(received)
		}
	}()

	// Read and send messages
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}

		// Prepend timestamp if not a /name command
		var message string
		if strings.HasPrefix(text, "/name ") {
			message = text
		} else {
			timestamp := time.Now().Format(time.RFC3339Nano)
			message = fmt.Sprintf("[LATENCY]|%s|%s", timestamp, text)
		}

		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Failed to send message:", err)
		}
	}
}
