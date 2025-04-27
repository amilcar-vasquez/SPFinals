package tcp

import (
    "bufio"
    "fmt"
    "net"
    "sync"
)

var clients = make(map[net.Conn]string)
var mutex = &sync.Mutex{}

func StartServer() {
    listener, err := net.Listen("tcp", ":9000")
    if err != nil {
        fmt.Println("Error starting TCP server:", err)
        return
    }
    defer listener.Close()

    fmt.Println("TCP Server started on :9000")

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Connection error:", err)
            continue
        }

        mutex.Lock()
        clients[conn] = conn.RemoteAddr().String()
        mutex.Unlock()

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer func() {
        mutex.Lock()
        delete(clients, conn)
        mutex.Unlock()
        conn.Close()
    }()

    reader := bufio.NewReader(conn)
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Client disconnected:", conn.RemoteAddr())
            return
        }
        broadcast(message, conn)
    }
}

func broadcast(message string, sender net.Conn) {
    mutex.Lock()
    defer mutex.Unlock()
    for conn := range clients {
        if conn != sender {
            fmt.Fprintf(conn, "%s", message)
        }
    }
}
