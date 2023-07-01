package mongoload

import (
	"sync"

	"github.com/montanaflynn/stats"
)

type Histogram interface {
	Len() int
	Min() (float64, error)
	Max() (float64, error)
	Mean() (float64, error)
	Percentile(float64) (float64, error)

	Update(float64)
}

type StandardHistogram struct {
	mutex sync.Mutex
	data  []float64
}

func NewHistogram() Histogram {
	data := make([]float64, 0)
	hist := &StandardHistogram{
		data: data,
	}

	return Histogram(hist)
}

func (h *StandardHistogram) Len() int {
	return len(h.data)
}

func (h *StandardHistogram) Min() (float64, error) {
	return stats.Min(h.data)
}

func (h *StandardHistogram) Max() (float64, error) {
	return stats.Max(h.data)
}

func (h *StandardHistogram) Mean() (float64, error) {
	return stats.Mean(h.data)
}

func (h *StandardHistogram) Percentile(percentile float64) (float64, error) {
	return stats.Percentile(h.data, percentile)
}

func (h *StandardHistogram) Update(interval float64) {
	h.mutex.Lock()
	h.data = append(h.data, interval)
	h.mutex.Unlock()
}
