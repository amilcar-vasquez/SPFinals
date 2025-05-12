package main

import (
	"net"
	"sync"
	"time"
)

type Client struct {
	Addr     *net.UDPAddr
	LastSeen time.Time
	Name     string
}

var (
	clients = make(map[string]Client)
	mu      sync.Mutex
)
