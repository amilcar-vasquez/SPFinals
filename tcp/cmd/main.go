//file: main.go

package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	// Define a flag to set the port number
	port := flag.Int("port", 4000, "Port to listen on")
	flag.Parse()

	// Start listening on the specified port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("Server listening on port %d\n", *port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}

		go loggingMiddleware(handleConnection)(conn)
	}
}
