// file: cmd/tcpHandlers.go
package main

import (
	"fmt"
	"github.com/fatih/color"
	"net"
	"os"
	"strings"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()

	// Ask nickname
	conn.Write([]byte("Please enter your nickname: "))
	nickBuf := make([]byte, 64)
	n, err := conn.Read(nickBuf)
	if err != nil {
		conn.Write([]byte("Failed to read nickname. Disconnecting.\n"))
		return
	}

	nickname := strings.TrimSpace(string(nickBuf[:n]))
	if nickname == "" {
		nickname = remoteAddr
	}

	manager.Add(conn)
	manager.SetNickname(conn, nickname)
	conn.Write([]byte(color.GreenString(fmt.Sprintf("Welcome to the chat, %s!\n", nickname))))
	manager.Broadcast(conn, color.YellowString(fmt.Sprintf("%s has joined the chat.", nickname)))

	// Logging
	fileName := strings.ReplaceAll(remoteAddr, ":", "_") + ".log"
	file, err := createLogFile(fileName)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	defer file.Close()

	// Metrics and timeout
	metrics.RegisterClient(conn)
	defer metrics.UnregisterClient(conn)
	reset := make(chan struct{})
	done := make(chan struct{})
	go StartTimer(conn, 600*time.Second, reset, done)

	// 📦 Extracted the heavy logic
	handleMessages(conn, nickname, file, reset, done)

	manager.Remove(conn)
}

func handleCommand(conn net.Conn, input string) bool {
	switch {
	case input == "/time":
		currentTime := time.Now().Format("15:04:05 MST 2006-01-02")
		conn.Write([]byte("Server time: " + currentTime + "\n"))
		return true

	case input == "/ping":
		metrics.RecordPingSent(conn)
		conn.Write([]byte("ping\n"))
		return true

	case input == "pong":
		metrics.RecordPingLatency(conn)
		conn.Write([]byte("pong received. Latency updated.\n"))
		return true

	case input == "/quit":
		conn.Write([]byte("Goodbye!\n"))
		conn.Close()
		return false

	case strings.HasPrefix(input, "/echo "):
		echoMessage := strings.TrimPrefix(input, "/echo ")
		conn.Write([]byte(echoMessage + "\n"))
		return true

	case strings.HasPrefix(input, "/nick "):
		newNick := strings.TrimSpace(strings.TrimPrefix(input, "/nick "))
		if newNick == "" {
			conn.Write([]byte("Nickname cannot be empty.\n"))
			return true
		}
		oldNick := manager.SetNickname(conn, newNick)
		if oldNick != newNick {
			conn.Write([]byte(fmt.Sprintf("Your nickname has been changed to %s\n", newNick)))
			manager.Broadcast(conn, fmt.Sprintf("%s changed nickname to %s", oldNick, newNick))
		} else {
			conn.Write([]byte("You already have that nickname.\n"))
		}
		return true
	case input == "/metrics":
		snapshot := metrics.GetSnapshot()
		conn.Write([]byte(snapshot + "\n"))
		return true

	case strings.HasPrefix(input, "/throughput"):
		addr := conn.RemoteAddr().String()
		m := metrics.ClientsData[addr]
		if m == nil {
			conn.Write([]byte("No data available.\n"))
		} else {
			duration := time.Since(m.ConnectionTime).Seconds()
			tp := float64(m.TotalBytes) / duration
			conn.Write([]byte(fmt.Sprintf("Throughput: %.2f bytes/sec\n", tp)))
		}
		return true

	default:
		conn.Write([]byte("Unknown command.\n"))
		return true
	}
}

func handleMessages(conn net.Conn, nickname string, file *os.File, reset, done chan struct{}) {
	buf := make([]byte, 1024)
	remoteAddr := conn.RemoteAddr().String()

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				fmt.Printf("Connection from %s was closed after timeout at %s\n", remoteAddr, time.Now().Format("03:04:05 PM"))
			} else {
				fmt.Printf("Unexpected read error from %s: %v\n", remoteAddr, err)
			}
			break
		}

		reset <- struct{}{}

		userInput := strings.TrimSpace(string(buf[:n]))
		if userInput == "" {
			continue
		}

		if strings.HasPrefix(userInput, "/") {
			if ok := handleCommand(conn, userInput); !ok {
				break
			}
			continue
		}

		truncatedMessage := truncateMessage(userInput, 100)

		_ = logMessage(file, "Client", truncatedMessage)
		broadcastMsg := fmt.Sprintf("%s: %s", nickname, truncatedMessage)
		_ = logMessage(file, "Broadcast", broadcastMsg)

		manager.Broadcast(conn, broadcastMsg)
		_ = logMessage(file, "Server", truncatedMessage)

		metrics.IncrementMessage(conn)
	}

	done <- struct{}{}
}
