package driver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kuzxnia/loadbot/lbot/config"
	"github.com/kuzxnia/loadbot/lbot/database"
	"github.com/kuzxnia/loadbot/lbot/schema"
	"github.com/samber/lo"
)

// todo: split this function to setup and to starting workers
type Worker struct {
	Metrics     *Metrics
	ctx         context.Context
	cfg         *config.Config
	job         *config.Job
	wg          sync.WaitGroup
	db          database.Client
	handler     JobHandler
	rateLimiter Limiter
	pool        JobPool
	dataPool    schema.DataPool
	ticker      *time.Ticker
	done        bool
}

func NewWorker(ctx context.Context, cfg *config.Config, job *config.Job, dataPool schema.DataPool) (*Worker, error) {
	// todo: check errors
	worker := new(Worker)
	worker.ctx = ctx
	worker.cfg = cfg
	worker.job = job
	worker.wg.Add(int(job.Connections))
	worker.pool = NewJobPool(job)
	worker.rateLimiter = NewLimiter(job)
	worker.Metrics = NewMetrics(job.Name)
	worker.done = false

	// introduce no db worker
	if job.Type != string(config.Sleep) {
		db, err := database.NewMongoClient(cfg.ConnectionString, job, job.GetSchema())
		if err != nil {
			return nil, err
		}
		worker.db = db
	}

	worker.dataPool = dataPool
	worker.handler = NewJobHandler(job, worker.db, dataPool)
	return worker, nil
}

func (w *Worker) Work() {
	fmt.Printf("Starting job: %s\n", lo.If(w.job.Name != "", w.job.Name).Else(w.job.Type))
	// something wrong with context propagation change this
	// go func() {
	// 	select {
	// 	case <-w.ctx.Done():
	// 		w.Cancel()
	// 	}
	// }()

	for i := 0; i < int(w.job.Connections); i++ {
		go func() {
			defer w.wg.Done()
			for w.pool.SpawnJob() {
				w.rateLimiter.Take()
				// perform operation

				w.Metrics.Meter(w.handler.Execute)

				w.pool.MarkJobDone()
			}
		}()
	}
	w.wg.Wait()
	w.done = true
  // todo: set end date
	// w.Report.SetDuration(time.Since(w.startTime))
}

func (w *Worker) InitMetrics() {
	w.Metrics.Init()
}

// todo: fix wrong place invalid
func (w *Worker) ExtendCopySavedFieldsToDataPool() {
	if w.dataPool != nil && (w.job.Type == string(config.Write) || w.job.Type == string(config.BulkWrite)) {
		w.dataPool.ExtendGeneratorMapperFields(schema.DefaultGeneratorFieldMapper)
	}
}

func (w *Worker) Summary() {
	// w.Report.Summary(nil)
}

func (w *Worker) Cancel() {
	fmt.Printf("Task canceled\n")
	w.pool.Cancel()
	w.Close()
}

func (w *Worker) IsDone() bool {
	return w.done
}

func (w *Worker) Close() {
	w.done = true
	if w.job.Type != string(config.Sleep) {
		w.db.Disconnect()
	}
	if w.ticker != nil {
		w.ticker.Stop()
	}
}

func (w *Worker) JobName() string {
	return w.job.Name
}

func (w *Worker) RequestedOperations() uint64 {
	return w.job.Operations
}

func (w *Worker) RequestedDurationSeconds() uint64 {
	return uint64(w.job.Duration.Seconds())
}
