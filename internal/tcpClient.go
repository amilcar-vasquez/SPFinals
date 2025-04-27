package internal

import (
    "net"
    "time"
)

type TCPClient struct {
    conn net.Conn
}

func NewTCPClient(address string) (*TCPClient, error) {
    conn, err := net.Dial("tcp", address)
    if err != nil {
        return nil, err
    }
    return &TCPClient{conn: conn}, nil
}

func (c *TCPClient) Send(msg string) error {
    c.conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
    _, err := c.conn.Write([]byte(msg + "\n"))
    return err
}

func (c *TCPClient) Receive() (string, error) {
    buf := make([]byte, 4096)
    c.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
    n, err := c.conn.Read(buf)
    if err != nil {
        return "", err
    }
    return string(buf[:n]), nil
}

func (c *TCPClient) Close() {
    c.conn.Close()
}
