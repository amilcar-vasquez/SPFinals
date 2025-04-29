// file: tcpHandlers.go

package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	manager.Add(conn)
	defer manager.Remove(conn)

	buf := make([]byte, 1024)
	remoteAddr := conn.RemoteAddr().String()
	fileName := strings.ReplaceAll(remoteAddr, ":", "_") + ".log"

	file, err := createLogFile(fileName)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	defer file.Close()

	timeout := 30 * time.Second
	reset := make(chan struct{})
	done := make(chan struct{})

	go StartTimer(conn, timeout, reset, done)

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

		if err := logMessage(file, "Client", truncatedMessage); err != nil {
			fmt.Println("Error logging client input:", err)
			break
		}

		client := manager.GetClient(conn)
		name := remoteAddr
		if client != nil && client.nickname != "" {
			name = client.nickname
		}

		broadcastMsg := fmt.Sprintf("%s: %s", name, truncatedMessage)

		if err := logMessage(file, "Broadcast", broadcastMsg); err != nil {
			fmt.Println("Error logging broadcast message:", err)
			break
		}

		manager.Broadcast(conn, broadcastMsg)

		conn.Write([]byte("You said: " + truncatedMessage + "\n"))

		response := personalityResponses(truncatedMessage)

		if err := logMessage(file, "Server", response); err != nil {
			fmt.Println("Error logging server response:", err)
			break
		}

		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			fmt.Println("Error writing to client:", err)
			break
		}
	}

	done <- struct{}{}
}

func handleCommand(conn net.Conn, input string) bool {
	switch {
	case input == "/time":
		currentTime := time.Now().Format("15:04:05 MST 2006-01-02")
		conn.Write([]byte("Server time: " + currentTime + "\n"))
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

	default:
		conn.Write([]byte("Unknown command.\n"))
		return true
	}
}
