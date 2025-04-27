package main

import (
    "github.com/amilcar-vasquez/SPFinal/internal"
    "net/http"
    "github.com/gorilla/websocket"
    "log"
    "sync"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

type ClientSession struct {
    conn     *websocket.Conn
    tcp      *internal.TCPClient
    udp      *internal.UDPClient
    protocol string // "tcp" or "udp"
    mu       sync.Mutex
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    wsConn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("WebSocket upgrade failed:", err)
        return
    }
    defer wsConn.Close()

    session := &ClientSession{
        conn:     wsConn,
        protocol: "tcp", // Default to TCP
    }

    session.tcp, _ = internal.NewTCPClient("127.0.0.1:9000") // TCP server addr
    session.udp, _ = internal.NewUDPClient("127.0.0.1:9001") // UDP server addr

    go session.listenFromServer()

    for {
        _, msg, err := wsConn.ReadMessage()
        if err != nil {
            internal.Logger.Println("Client disconnected:", err)
            break
        }

        fullMsg := string(msg)

        if fullMsg == "switch:tcp" {
            session.mu.Lock()
            session.protocol = "tcp"
            session.mu.Unlock()
            continue
        } else if fullMsg == "switch:udp" {
            session.mu.Lock()
            session.protocol = "udp"
            session.mu.Unlock()
            continue
        }

        session.mu.Lock()
        if session.protocol == "tcp" {
            session.tcp.Send(fullMsg)
        } else {
            session.udp.Send(fullMsg)
        }
        session.mu.Unlock()

        internal.Logger.Println("Sent:", fullMsg)
    }

    session.tcp.Close()
    session.udp.Close()
}

func (s *ClientSession) listenFromServer() {
    for {
        var msg string
        var err error

        s.mu.Lock()
        if s.protocol == "tcp" {
            msg, err = s.tcp.Receive()
        } else {
            msg, err = s.udp.Receive()
        }
        s.mu.Unlock()

        if err != nil {
            continue // In production maybe handle this better
        }

        s.conn.WriteMessage(websocket.TextMessage, []byte(msg))
    }
}

func main() {
    internal.InitLogger()
    defer internal.CloseLogger()

    r := http.NewServeMux()
    r.Handle("/", http.FileServer(http.Dir("./ui/html")))
    r.HandleFunc("/ws", wsHandler)

    log.Println("Starting server at :4000")
    log.Fatal(http.ListenAndServe(":4000", r))

    log.Println("Starting server at :4000")
    log.Fatal(http.ListenAndServe(":4000", nil))
}
