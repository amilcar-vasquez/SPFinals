# 📡 UDP Chat Server in Go

This project is a **concurrent chat server** implemented using Go's `net` package over **UDP**. It supports multiple clients, message broadcasting, inactivity-based disconnection, real-time metrics logging, and simulated networking conditions.

---

## 🚀 Features

✅ Clients can send and receive messages  
✅ Messages are broadcast to all connected clients  
✅ Handles name setting and graceful client disconnection  
✅ Uses goroutines and synchronization (`sync.Mutex`) for safe concurrency  
✅ Tracks metrics: message count, broadcasts, bytes sent, active clients, uptime  
✅ Simulates real-world network conditions: latency, packet loss  
✅ Logs system performance to `metrics.log`

---

## 📁 Project Structure

├── main.go # Starts metrics logger and server
├── server.go # Core UDP server logic
├── broadcast.go # Broadcasts messages to all clients
├── client.go # Client struct and global map
├── cleanup.go # Periodic cleanup of inactive clients
├── metrics.go # Metrics tracking and logging
├── log_metrics.go # Logs metrics to file and console
├── test-client.go # Example UDP client (CLI)
└── metrics.log # (Generated) log file for stats



---

## ⚙️ How to Run

### 🖥️ Run the Server

```bash
go run .

```

## You'll see:
UDP server started on port 9001
🧹 Cleaning up inactive clients...


## Open two other terminals:
cd test
```bash
go run test-client.go
```

## You'll see:
Connected to UDP server at 127.0.0.1:9001
Type your messages below. Use '/name YourName' to set your name.

## Set a name:
/name Darwin

## Send a message:
Hello everyone!

## 📊 Metrics Logging
Metrics are recorded every 15 seconds in metrics.log and include:

Uptime

Total messages

Total broadcasts

Bytes sent

Active clients

## EXAMPLE:
[2025-05-11 14:05:10] Uptime: 45s | Messages: 6 | Broadcasts: 5 | Bytes: 480 | Active Clients: 2

## 🧬 Simulating Network Conditions
Your server randomly drops 10% of outgoing messages to simulate packet loss.

## 🔁 Auto Cleanup
Inactive clients (no message in 1 minute) are automatically removed.

Logged every 10 seconds:
🧹 Cleaning up inactive clients...
Removing inactive client: 127.0.0.1:54321


## 👨‍💻 Author
Darwin Ramos
Systems Final Project – UDP Chat Server (Go)