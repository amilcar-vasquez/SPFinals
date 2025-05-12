// file: logger.go
package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func loggingMiddleware(next func(net.Conn)) func(net.Conn) {
	return func(conn net.Conn) {
		remoteAddr := conn.RemoteAddr().String()
		// Log the connection details including address and connection time
		fmt.Printf("New connection from %s at %s\n", remoteAddr, time.Now().Format("03:04:05 PM"))
		// Log connection close time
		defer func() {
			fmt.Printf("Connection from %s closed at %s\n", remoteAddr, time.Now().Format("03:04:05 PM"))
		}()
		next(conn)
	}
}

// helper function to log messages
func logMessage(file *os.File, prefix, message string) error {
	timestamp := time.Now().Format("03:04:05 PM")
	_, err := file.WriteString(fmt.Sprintf("[%s] %s: %s\n", timestamp, prefix, message))
	return err
}

/* createLogFile creates a log file with the given name.
It opens the file in append mode and creates it if it doesn't exist.
The file permissions are set to 0644. */

func createLogFile(fileName string) (*os.File, error) {
	// Ensure the "log" folder exists
	if err := os.MkdirAll("log", 0755); err != nil {
		return nil, err
	}
	// Open the file in the "log" folder in append mode, create it if it doesn't exist
	return os.OpenFile("log/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}
