// metrics.go
package main

import (
	"sync"
	"time"
)

var (
	// Existing
	startTime        = time.Now()
	totalMessages    int
	totalBroadcasts  int
	totalBytes       int
	metricsMutex     sync.Mutex

	// New for latency and packet loss
	totalLatency     time.Duration
	latencySamples   int
	simulatedDrops   int
)


func incrementMessageCount(size int) {
	metricsMutex.Lock()
	totalMessages++
	totalBytes += size
	metricsMutex.Unlock()
}

func incrementBroadcastCount() {
	metricsMutex.Lock()
	totalBroadcasts++
	metricsMutex.Unlock()
}

func getMetricsSnapshot() (int, int, int, int, time.Duration, time.Duration, int, int) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	avgLatency := time.Duration(0)
	if latencySamples > 0 {
		avgLatency = totalLatency / time.Duration(latencySamples)
	}

	activeClients := len(clients)
	return totalMessages, totalBroadcasts, totalBytes, activeClients, time.Since(startTime), avgLatency, latencySamples, simulatedDrops
}


func addLatencySample(d time.Duration) {
	metricsMutex.Lock()
	totalLatency += d
	latencySamples++
	metricsMutex.Unlock()
}

func incrementSimulatedDrops() {
	metricsMutex.Lock()
	simulatedDrops++
	metricsMutex.Unlock()
}
