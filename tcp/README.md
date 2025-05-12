# 🧠 TCP Chat Server in Go

A robust, concurrent TCP chat server built with Go, featuring real-time messaging, client management, and performance metrics.

## 🚀 Features

* **Real-Time Messaging**: Clients can send and receive messages instantly.
* **Broadcast System**: Messages are broadcasted to all connected clients.
* **Graceful Disconnection**: Handles client disconnections smoothly.
* **Concurrency**: Utilizes goroutines and mutexes for efficient concurrent handling.
* **Performance Metrics**: Tracks latency, throughput, and packet loss.
* **Command Support**: Supports various client commands for enhanced interaction.


### Running the Server

```bash
make run
```

The server listens on the specified port (default: 4000).

### Connecting Clients

Clients can connect using `netcat` or any TCP client:([Go Packages][1])

```bash
nc localhost 4000
```



## 💬 Client Commands

* `/time` – Displays the current server time.
* `/ping` – Measures latency between client and server.
* `/nick <name>` – Changes the client's nickname.
* `/echo <message>` – Echoes the message back to the client.
* `/metrics` – Displays server performance metrics.
* `/throughput` – Shows the client's data throughput.
* `/quit` – Disconnects the client from the server.

## 📈 Performance Metrics

The server tracks and reports:

* **Total Messages**: Number of messages processed.
* **Active Clients**: Current connected clients.
* **Client-Specific Data**:

  * Messages sent
  * Connection duration
  * Data throughput
  * Packet loss
  * Latency samples([Medium][2])

Use the `/metrics` command to view these metrics.

## 🧪 Testing & Optimization

### Simulating Network Conditions

Use tools like `tc` on Linux to simulate latency and packet loss:

```bash
sudo tc qdisc add dev lo root netem delay 100ms loss 10%
```

### Handling Edge Cases

* **Sudden Disconnects**: The server detects and handles unexpected client disconnections.
* **High Message Volume**: Test the server's performance under heavy load using tools like `ab` or custom scripts.

## 🧹 Code Structure

* `main.go` – Entry point; starts the server and listens for connections.
* `cmd/clientHandlers.go` – Manages client connections and broadcasting.
* `cmd/tcpHandlers.go` – Handles client commands and message processing.
* `metrics.go` – Collects and reports performance metrics.
* `logger.go` – Logs server events and client interactions.
* `overflowProtector.go` – Ensures messages do not exceed maximum length.
* `timer.go` – Implements inactivity timeout for clients.

## Demo Video URL:


