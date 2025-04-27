package udp

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

func StartClient() {
    serverAddr, err := net.ResolveUDPAddr("udp", "localhost:9001")
    if err != nil {
        fmt.Println("Address resolve error:", err)
        return
    }

    conn, err := net.DialUDP("udp", nil, serverAddr)
    if err != nil {
        fmt.Println("Connection error:", err)
        return
    }
    defer conn.Close()

    go listenServer(conn)

    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := scanner.Text()
        conn.Write([]byte(text))
    }
}

func listenServer(conn *net.UDPConn) {
    buf := make([]byte, 1024)
    for {
        n, _, err := conn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println("Disconnected from server.")
            os.Exit(0)
        }
        fmt.Println(">> " + string(buf[:n]))
    }
}
