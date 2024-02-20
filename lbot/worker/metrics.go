package worker

import (
	"fmt"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/google/uuid"
	"github.com/kuzxnia/loadbot/lbot/config"
)

type Metrics struct {
	requests        *metrics.Counter
	requestsError   *metrics.Counter
	requestDuration *metrics.Summary
	startTime       time.Time
	// ResponseSize    *metrics.Histogram
}

func NewMetrics(job *config.Job) *Metrics {
	jobLabel := fmt.Sprintf(`{job="%s",job_uuid="%s",job_type="%s"}`, job.Name, uuid.New().String(), job.Type)

	return &Metrics{
		requests:        metrics.NewCounter("requests_total" + jobLabel),
		requestsError:   metrics.NewCounter("requests_error" + jobLabel),
		requestDuration: metrics.NewSummary("requests_duration_seconds" + jobLabel),
		// ResponseSize:    metrics.NewHistogram("requests_size"),
	}
}

func (m *Metrics) Init() {
	m.startTime = time.Now()
}

func (m *Metrics) Meter(handler func() error) {
	startTime := time.Now()

	error := handler()

	// todo: handle size
	m.requestDuration.UpdateDuration(startTime)
	m.requests.Inc()
	if error != nil {
		m.requestsError.Inc()
	}
}

func (m *Metrics) Rps() uint64 {
	duration := time.Since(m.startTime).Seconds()
	if duration == 0 {
		return 0
	}
	return uint64(float64(m.requests.Get()) / duration)
}

func (m *Metrics) Requests() uint64 {
	return m.requests.Get()
}

func (m *Metrics) ErrorRate() float32 {
	return float32(m.requestsError.Get()) / float32(m.requests.Get())
}

func (m *Metrics) DurationSeconds() uint64 {
	return uint64(time.Since(m.startTime).Round(time.Second).Seconds())
}
