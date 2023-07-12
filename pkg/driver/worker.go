package driver

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kuzxnia/mongoload/pkg/config"
	"github.com/kuzxnia/mongoload/pkg/database"
	"github.com/kuzxnia/mongoload/pkg/rps"
)

type Worker interface {
	Work()
	ExecuteJob()
	Cancel()
	Summary()
}

func NewWorker(cfg *config.Job) (Worker, error) {
	switch cfg.Type {
	case string(config.Write):
		worker := new(InsertWorker)
		worker.wg.Add(int(cfg.Connections))
		worker.Statistic = NewWriteStats()
		worker.pool = NewJobPool(cfg)
		worker.rateLimiter = rps.NewLimiter(cfg)
		worker.startTime = time.Now()
    // todo: init db

		return Worker(worker), nil
	default:
		return nil, errors.New("Invalid job type")
	}
}

type BaseWorker struct {
	wg          sync.WaitGroup
	db          database.Client
	rateLimiter rps.Limiter
	pool        JobPool
	startTime   time.Time
}

type InsertWorker struct {
	BaseWorker
	Statistic Stats
}

func (w *InsertWorker) Work() {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	go func() {
		<-interruptChan
		w.Cancel()
	}()
}

func (w *InsertWorker) ExecuteJob() {
	// start       time.Time

	defer w.wg.Done()

	for w.pool.SpawnJob() {
		w.rateLimiter.Take()
		// perform operation
		start := time.Now()
		// do sth with is error
		_, error := w.db.InsertOneOrMany()
		elapsed := time.Since(start)
		w.Statistic.Add(float64(elapsed.Milliseconds()), error)

		// add debug of some kind
		if error != nil {
			// todo: debug
			log.Debug(error)
		}

		w.pool.MarkJobDone()
	}
}

func (w *InsertWorker) Summary() {
	elapsed := time.Since(w.startTime)
	requestsDone := w.pool.GetRequestsDone()
	rps := float64(requestsDone) / elapsed.Seconds()

	fmt.Printf("\nTime took %f s\n", elapsed.Seconds())
	fmt.Printf("Total operations: %d\n", requestsDone)
  w.Statistic.Summary()
	fmt.Printf("Requests per second: %f rp/s\n", rps)
	if batch := w.db.GetBatchSize(); batch > 1 {
		fmt.Printf("Operations per second: %f op/s\n", float64(requestsDone*batch)/elapsed.Seconds())
	}
}

func (w *InsertWorker) Cancel() {
	w.pool.Cancel()
}

type JobPool interface {
	SpawnJob() bool
	MarkJobDone()
	Cancel()
	GetRequestsStarted() uint64
	GetRequestsDone() uint64
}

func NewJobPool(cfg *config.Job) JobPool {
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
