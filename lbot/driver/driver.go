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
	ctx         context.Context
	cfg         *config.Config
	job         *config.Job
	wg          sync.WaitGroup
	db          database.Client
	handler     JobHandler
	rateLimiter Limiter
	pool        JobPool
	dataPool    schema.DataPool
	metrics     *Metric
	// Report      Report
	ticker *time.Ticker
}

func NewWorker(ctx context.Context, cfg *config.Config, job *config.Job, dataPool schema.DataPool, metric *Metric) (*Worker, error) {
	// todo: check errors
	fmt.Printf("Starting job: %s\n", lo.If(job.Name != "", job.Name).Else(job.Type))
	worker := new(Worker)
	worker.ctx = ctx
	worker.cfg = cfg
	worker.job = job
	worker.wg.Add(int(job.Connections))
	// worker.Report = NewReport(job)
	worker.metrics = metric
	worker.pool = NewJobPool(job)
	worker.rateLimiter = NewLimiter(job)

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

				w.metrics.Meter(w.handler.Execute)

				w.pool.MarkJobDone()
			}
		}()
	}
	w.wg.Wait()
	// w.Report.SetDuration(time.Since(w.startTime))
}

func (w *Worker) InitIntervalReportingSummary(logs chan string) {
	reportingFormat := w.job.GetReport()
	if reportingFormat == nil || reportingFormat.Interval == 0 {
		// log.Info("Interval reporting skipped")
		return
	}

	w.ticker = time.NewTicker(reportingFormat.Interval)
	go func(worker *Worker) {
		for range w.ticker.C {
			// todo: get metrics and push
			// worker.Report.SetDuration(time.Since(w.startTime))
			// worker.Report.Summary(logs)
		}
	}(w)
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

func (w *Worker) Close() {
	if w.job.Type != string(config.Sleep) {
		w.db.Disconnect()
	}
	if w.ticker != nil {
		w.ticker.Stop()
	}
}
