package worker

import (
	"sync"
	"sync/atomic"
	"time"
)

type JobPool interface {
	SpawnJob() bool
	MarkJobDone()

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
		requestsStarted := atomic.AddUint64(&w.requestsStarted, 1)
		return requestsStarted <= w.requestsNumber
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

func (w *deductionJobPool) GetRequestsDone() uint64 {
	return atomic.LoadUint64(&w.requestsDone)
}

type timerJobPool struct {
	start        time.Time
	duration     time.Duration
	requestsDone uint64

	done  chan struct{}
	close sync.Once
}

func NewTimerJobPool(duration time.Duration) JobPool {
	if duration < 0 {
		panic("duration must be positive")
	}

	pool := &timerJobPool{
		start:        time.Now(),
		duration:     duration,
		requestsDone: 0,
		done:         make(chan struct{}),
	}
	go func() {
		time.AfterFunc(duration, func() {
			pool.Cancel()
		})
	}()
	return JobPool(pool)
}

func (w *timerJobPool) SpawnJob() bool {
	select {
	case <-w.done:
		return false
	default:
		return true
	}
}

func (w *timerJobPool) MarkJobDone() {
	atomic.AddUint64(&w.requestsDone, 1)
}

func (w *timerJobPool) Cancel() {
	w.close.Do(func() { close(w.done) })
}

func (w *timerJobPool) GetRequestsDone() uint64 {
	return atomic.LoadUint64(&w.requestsDone)
}
