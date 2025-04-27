package metrics

import (
    "time"
)

type Metrics struct {
    SentMessages     int
    ReceivedMessages int
    TotalBytes       int
    StartTime        time.Time
}

func (m *Metrics) Start() {
    m.StartTime = time.Now()
}

func (m *Metrics) RecordSend(bytes int) {
    m.SentMessages++
    m.TotalBytes += bytes
}

func (m *Metrics) RecordReceive(bytes int) {
    m.ReceivedMessages++
    m.TotalBytes += bytes
}

func (m *Metrics) Throughput() float64 {
    duration := time.Since(m.StartTime).Seconds()
    if duration == 0 {
        return 0
    }
    return float64(m.TotalBytes) / duration // bytes per second
}
