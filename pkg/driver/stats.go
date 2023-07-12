package driver

import (
	"fmt"
	"sync"

	"github.com/kuzxnia/mongoload/pkg/config"
	"github.com/montanaflynn/stats"
	"go.mongodb.org/mongo-driver/mongo"
)

type Stats interface {
	Len() int
	Min() (float64, error)
	Max() (float64, error)
	Mean() (float64, error)
	Percentile(float64) (float64, error)

	Add(float64, error)

	Summary()
}

func NewStatistics(job *config.Job) Stats {
	switch job.Type {
	case string(config.Write):
	case string(config.BulkWrite):
		return Stats(
			&WriteStats{
				BaseStats: BaseStats{
					data: make([]float64, 0),
				},
			},
		)
	case string(config.Read):
		return Stats(
			&ReadStats{
				BaseStats: BaseStats{
					data: make([]float64, 0),
				},
			},
		)
	case string(config.Update):
		return Stats(
			&UpdateStats{
				BaseStats: BaseStats{
					data: make([]float64, 0),
				},
			},
		)
	default:
		// todo change
		panic("Invalid job type")
	}
	return nil
}

type BaseStats struct {
	mutex         sync.RWMutex
	data          []float64
	timeoutErrors uint64
	otherErrors   uint64
}

func (s *BaseStats) Len() int               { return len(s.data) }
func (s *BaseStats) Min() (float64, error)  { return stats.Min(s.data) }
func (s *BaseStats) Max() (float64, error)  { return stats.Max(s.data) }
func (s *BaseStats) Mean() (float64, error) { return stats.Mean(s.data) }
func (s *BaseStats) Percentile(percentile float64) (float64, error) {
	return stats.Percentile(s.data, percentile)
}

type ReadStats struct {
	BaseStats
	noDocumentsFoundError uint64
}

func (s *ReadStats) Add(interval float64, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data = append(s.data, interval)
	if err != nil {
		if mongo.IsTimeout(err) {
			s.timeoutErrors++
		} else if err == mongo.ErrNoDocuments {
			s.noDocumentsFoundError++
		} else {
			s.otherErrors++
		}
	}
}

func (s *ReadStats) Summary() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if len := s.Len(); len != 0 {
		errors := int(s.timeoutErrors + s.otherErrors + s.noDocumentsFoundError)
		wmean, _ := s.Mean()
		wp50, _ := s.Percentile(50)
		wp90, _ := s.Percentile(90)
		wp99, _ := s.Percentile(99)
		fmt.Printf(
			"Total read ops: %d, successful: %d, errors: (timeout: %d, noDataFound: %d, other: %d), error rate: %.2f%% \n",
			len, len-errors, s.timeoutErrors, s.noDocumentsFoundError, s.otherErrors, float64(errors)/float64(len)*100,
		)
		fmt.Printf("Read AVG: %.2fms, P50: %.2fms, P90: %.2fms P99: %.2fms\n", wmean, wp50, wp90, wp99)
	}
}

type WriteStats struct {
	BaseStats
}

func (s *WriteStats) Add(interval float64, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data = append(s.data, interval)
	if err != nil {
		if mongo.IsTimeout(err) {
			s.timeoutErrors++
		} else {
			s.otherErrors++
		}
	}
}

func (s *WriteStats) Summary() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if len := s.Len(); len != 0 {
		errors := int(s.timeoutErrors + s.otherErrors)
		wmean, _ := s.Mean()
		wp50, _ := s.Percentile(50)
		wp90, _ := s.Percentile(90)
		wp99, _ := s.Percentile(99)
		fmt.Printf(
			"Total write ops: %d, successful: %d, errors: (timeout: %d, other: %d), error rate: %.2f%% \n",
			len, len-errors, s.timeoutErrors, s.otherErrors, float64(errors)/float64(len)*100,
		)
		fmt.Printf("Write AVG: %.2fms, P50: %.2fms, P90: %.2fms P99: %.2fms\n", wmean, wp50, wp90, wp99)
	}
}

type UpdateStats struct {
	BaseStats
}

func (s *UpdateStats) Add(interval float64, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data = append(s.data, interval)
	if err != nil {
		if mongo.IsTimeout(err) {
			s.timeoutErrors++
		} else {
			s.otherErrors++
		}
	}
}

func (s *UpdateStats) Summary() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if len := s.Len(); len != 0 {
		errors := int(s.timeoutErrors + s.otherErrors)
		wmean, _ := s.Mean()
		wp50, _ := s.Percentile(50)
		wp90, _ := s.Percentile(90)
		wp99, _ := s.Percentile(99)
		fmt.Printf(
			"Total Update ops: %d, successful: %d, errors: (timeout: %d, other: %d), error rate: %.2f%% \n",
			len, len-errors, s.timeoutErrors, s.otherErrors, float64(errors)/float64(len)*100,
		)
		fmt.Printf("Update AVG: %.2fms, P50: %.2fms, P90: %.2fms P99: %.2fms\n", wmean, wp50, wp90, wp99)
	}
}

// var ErrNoDocuments = errors.New("mongo: no documents in result")
// type ActorHistogram struct {
// 	Stats
// 	datach chan float64
// }

// func NewActorHistogram() Histogram {
// 	data := make([]float64, 0)
// 	hist := &ActorHistogram{
// 		BaseHistogram: BaseHistogram{
// 			data: data,
// 		},
//     datach: make(chan float64),
// 	}

// 	go hist.loop()

// 	return Histogram(hist)
// }

// func (h *ActorHistogram) Update(interval float64) {
// 	h.datach <- interval
// }

// func (h *ActorHistogram) loop() {
// 	for interval := range h.datach {
// 		h.data = append(h.data, interval)
// 	}
// }
