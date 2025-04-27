//file: main.go

package main

import (
    "flag"
    "fmt"
    "log"
	"github.com/amilcar-vasquez/SPFinal/tcp"
	"github.com/amilcar-vasquez/SPFinal/udp"
)

func main() {
    protocol := flag.String("protocol", "tcp", "Protocol to use: tcp or udp")
    flag.Parse()

    if *protocol == "tcp" {
        fmt.Println("Starting TCP Server...")
        tcp.StartServer()
    } else if *protocol == "udp" {
        fmt.Println("Starting UDP Server...")
        udp.StartServer()
    } else {
        log.Fatal("Unknown protocol. Use --protocol=tcp or --protocol=udp")
    }
}
