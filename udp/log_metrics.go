// log_metrics.go
package main

import (
	"fmt"
	"os"
	"time"
)

func StartMetricsLogger() {
	logFile, err := os.OpenFile("metrics.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("❌ Error opening metrics log file:", err)
		return
	}
	defer logFile.Close()

	for {
		time.Sleep(15 * time.Second)
		msgs, bcasts, bytes, active, uptime, avgLatency, latencySamples, drops := getMetricsSnapshot()


		timestamp := time.Now().Format("2006-01-02 15:04:05")
		log := fmt.Sprintf(
			"[%s] Uptime: %v | Messages: %d | Broadcasts: %d | Bytes: %d | Active Clients: %d | Avg Latency: %v | Samples: %d | Simulated Drops: %d\n",
			timestamp,
			uptime.Truncate(time.Second),
			msgs,
			bcasts,
			bytes,
			active,
			avgLatency,
			latencySamples,
			drops,
		)
		

		// Print to console
		//fmt.Println("\n📊 Metrics:")
		//fmt.Printf("%s", log)
		//fmt.Println("------")

		// Write to metrics.log
		if _, err := logFile.WriteString(log); err != nil {
			fmt.Println("❌ Error writing to metrics log:", err)
		}
	}
}
