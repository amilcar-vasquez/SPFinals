package tcp

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

func StartClient() {
    conn, err := net.Dial("tcp", "localhost:9000")
    if err != nil {
        fmt.Println("Error connecting:", err)
        return
    }
    defer conn.Close()

    go listenServer(conn)

    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := scanner.Text() + "\n"
        conn.Write([]byte(text))
    }
}

func listenServer(conn net.Conn) {
    reader := bufio.NewReader(conn)
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Disconnected from server.")
            os.Exit(0)
        }
        fmt.Print(">> " + message)
    }
}
