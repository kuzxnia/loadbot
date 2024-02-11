package driver

import (
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/samber/lo"
)

type Metrics struct {
	requests        *metrics.Counter
	requestsError   *metrics.Counter
	requestDuration *metrics.Summary
	startTime       time.Time
	// ResponseSize    *metrics.Histogram
}

func NewMetrics(job_name string) *Metrics {
	jobLabel := "{job=\"" + job_name + "\"}"

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

	// handle size
	m.requestDuration.UpdateDuration(startTime)
	m.requests.Inc()
	if error != nil {
		m.requestsError.Inc()
	}
}

func (m *Metrics) Rps() float32 {
	duration := time.Since(m.startTime)
	return lo.If(duration != 0, float32(m.requests.Get())/m.DurationSeconds()).Else(0)
}

func (m *Metrics) Requests() uint64 {
	return m.requests.Get()
}

func (m *Metrics) ErrorRate() float32 {
	return float32(m.requestsError.Get()) / float32(m.requests.Get())
}

func (m *Metrics) DurationSeconds() float32 {
	return float32(time.Since(m.startTime).Seconds())
}
