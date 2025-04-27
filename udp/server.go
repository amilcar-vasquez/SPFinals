package udp

import (
    "fmt"
    "net"
)

var clients = make(map[string]*net.UDPAddr)

func StartServer() {
    addr, err := net.ResolveUDPAddr("udp", ":9001")
    if err != nil {
        fmt.Println("Error resolving address:", err)
        return
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Println("Error starting UDP server:", err)
        return
    }
    defer conn.Close()

    fmt.Println("UDP Server started on :9001")

    buf := make([]byte, 1024)

    for {
        n, clientAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println("Read error:", err)
            continue
        }

        message := string(buf[:n])

        // Save client if new
        clients[clientAddr.String()] = clientAddr

        broadcast(conn, message, clientAddr)
    }
}

func broadcast(conn *net.UDPConn, message string, sender *net.UDPAddr) {
    for addrStr, addr := range clients {
        if addrStr != sender.String() {
            conn.WriteToUDP([]byte(message), addr)
        }
    }
}
