// file: timer.go
package main

import (
	"fmt"
	"net"
	"time"
)

func StartTimer(conn net.Conn, timeout time.Duration, reset <-chan struct{}, done <-chan struct{}) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-reset:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(timeout)
		case <-done:
			return // connection is finished, stop the timer
		case <-timer.C:
			fmt.Printf("Connection from %s timed out due to inactivity.\n", conn.RemoteAddr())
			conn.Close()
			return
		}
	}
}
