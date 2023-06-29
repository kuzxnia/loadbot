package worker

import (
	"sync"
	"sync/atomic"
)

type JobPool interface {
	SpawnJob() bool
	MarkJobDone()

	IsWorkerComplete() float64
	Cancel()

  GetRequestsDone() uint64
}

type deductionJobPool struct {
	requestsStarted uint64
	requestsDone    uint64
	requestsNumber  uint64

	done  chan struct{}
	close sync.Once
}

func NewDeductionJobPool(requestsNumber uint64) JobPool {
  pool := &deductionJobPool{
		requestsStarted: 0,
		requestsDone:    0,
		requestsNumber:  requestsNumber,
		done:            make(chan struct{}),
	}
  return JobPool(pool)
}

func (w *deductionJobPool) SpawnJob() bool {
	select {
	case <-w.done:
		return false

	default:
		requestsDone := atomic.AddUint64(&w.requestsDone, 1)
		return requestsDone <= w.requestsNumber
	}
}

func (w *deductionJobPool) MarkJobDone() {
	requestsDone := atomic.AddUint64(&w.requestsDone, 1)
	if requestsDone == w.requestsNumber {
		w.close.Do(func() { close(w.done) })
	}
}

func (w *deductionJobPool) Cancel() {
	w.close.Do(func() { close(w.done) })
}

func (w *deductionJobPool) IsWorkerComplete() float64 {
	select {
	case <-w.done:
		return 1.0

	default:
		requestsDone := atomic.AddUint64(&w.requestsDone, 1)
		return float64(requestsDone) / float64(w.requestsNumber)
	}
}

func (w *deductionJobPool) GetRequestsDone() uint64 {
  return atomic.LoadUint64(&w.requestsDone)
}
