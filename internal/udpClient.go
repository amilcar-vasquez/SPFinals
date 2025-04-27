package internal

import (
    "net"
    "time"
)

type UDPClient struct {
    conn   net.Conn
}

func NewUDPClient(address string) (*UDPClient, error) {
    conn, err := net.Dial("udp", address)
    if err != nil {
        return nil, err
    }
    return &UDPClient{conn: conn}, nil
}

func (c *UDPClient) Send(msg string) error {
    c.conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
    _, err := c.conn.Write([]byte(msg))
    return err
}

func (c *UDPClient) Receive() (string, error) {
    buf := make([]byte, 4096)
    c.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
    n, err := c.conn.Read(buf)
    if err != nil {
        return "", err
    }
    return string(buf[:n]), nil
}

func (c *UDPClient) Close() {
    c.conn.Close()
}
