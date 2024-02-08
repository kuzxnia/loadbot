package driver

import (
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/samber/lo"
)

var Stats = NewMetrics()

type Metric struct {
	RequestsTotal   *metrics.Counter
	RequestsError   *metrics.Counter
	RequestDuration *metrics.Summary
	startTime       time.Time
	// ResponseSize    *metrics.Histogram
}

func NewMetrics() *Metric {
	return &Metric{
		RequestsTotal:   metrics.NewCounter("requests_total"),
		RequestsError:   metrics.NewCounter("requests_error"),
		RequestDuration: metrics.NewSummary("requests_duration_seconds"),
		// ResponseSize:    metrics.NewHistogram("requests_size"),
	}
}

func (m *Metric) Init() {
	m.startTime = time.Now()
}

func (m *Metric) Meter(handler func() error) {
	startTime := time.Now()

	error := handler()

	// handle size
	m.RequestDuration.UpdateDuration(startTime)
	m.RequestsTotal.Inc()
	if error != nil {
		m.RequestsError.Inc()
	}
}

func (m *Metric) Rps() float32 {
	duration := time.Since(m.startTime)
	return lo.If(duration != 0, float32(m.RequestsTotal.Get())/float32(duration.Seconds())).Else(0)
}

func (m *Metric) Total() uint64 {
	return m.RequestsTotal.Get()
}

func (m *Metric) ErrorRate() float32 {
	return float32(m.RequestsError.Get()) / float32(m.RequestsTotal.Get())
}

// startTime := time.Now()

// error := w.handler.Handle()
// // w.Report.Add(duration, error)
// RequestDuration.UpdateDuration(startTime)
// // add debug of some kind

// RequestsTotal.Inc()
