package driver

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/kuzxnia/mongoload/pkg/config"
)

type JobPool interface {
	SpawnJob() bool
	MarkJobDone()
	Cancel()
	GetRequestsStarted() uint64
	GetRequestsDone() uint64
}

func NewJobPool(cfg *config.Job) JobPool {
	// todo: refactor this, add tracing
	if cfg.Duration == 0 && cfg.Operations == 0 {
		return JobPool(NewNoLimitTimerJobPool())
	} else if cfg.Duration != 0 {
		return JobPool(NewTimerJobPool(cfg.Duration))
	} else {
		return JobPool(NewDeductionJobPool(cfg.Operations))
	}
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

func (w *deductionJobPool) GetRequestsStarted() uint64 {
	return atomic.LoadUint64(&w.requestsStarted)
}

func (w *deductionJobPool) GetRequestsDone() uint64 {
	return atomic.LoadUint64(&w.requestsDone)
}

type timerJobPool struct {
	duration        time.Duration
	requestsStarted uint64
	requestsDone    uint64

	done  chan struct{}
	close sync.Once
}

func NewTimerJobPool(duration time.Duration) JobPool {
	if duration < 0 {
		panic("duration must be positive")
	}

	pool := &timerJobPool{
		requestsStarted: 0,
		requestsDone:    0,
		duration:        duration,
		done:            make(chan struct{}),
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
		atomic.AddUint64(&w.requestsStarted, 1)
		return true
	}
}

func (w *timerJobPool) MarkJobDone() {
	atomic.AddUint64(&w.requestsDone, 1)
}

func (w *timerJobPool) Cancel() {
	w.close.Do(func() { close(w.done) })
}

func (w *timerJobPool) GetRequestsStarted() uint64 {
	return atomic.LoadUint64(&w.requestsStarted)
}

func (w *timerJobPool) GetRequestsDone() uint64 {
	return atomic.LoadUint64(&w.requestsDone)
}

type noLimitTimerJobPool struct {
	requestsStarted uint64
	requestsDone    uint64

	done  chan struct{}
	close sync.Once
}

func NewNoLimitTimerJobPool() JobPool {
	pool := &noLimitTimerJobPool{
		requestsDone: 0,
		done:         make(chan struct{}),
	}
	return JobPool(pool)
}

func (w *noLimitTimerJobPool) SpawnJob() bool {
	select {
	case <-w.done:
		return false
	default:
		atomic.AddUint64(&w.requestsStarted, 1)
		return true
	}
}

func (w *noLimitTimerJobPool) MarkJobDone() {
	atomic.AddUint64(&w.requestsDone, 1)
}

func (w *noLimitTimerJobPool) Cancel() {
	w.close.Do(func() { close(w.done) })
}

func (w *noLimitTimerJobPool) GetRequestsStarted() uint64 {
	return atomic.LoadUint64(&w.requestsStarted)
}

func (w *noLimitTimerJobPool) GetRequestsDone() uint64 {
	return atomic.LoadUint64(&w.requestsDone)
}
