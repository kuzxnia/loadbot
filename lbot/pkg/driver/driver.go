package driver

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/kuzxnia/loadbot/lbot/pkg/config"
	"github.com/kuzxnia/loadbot/lbot/pkg/database"
	"github.com/kuzxnia/loadbot/lbot/pkg/logger"
	"github.com/kuzxnia/loadbot/lbot/pkg/schema"
	"github.com/samber/lo"
)

var log = logger.Default()


// todo: split this function to setup and to starting workers

func Torment(config *config.Config) {
	// todo: ping db, before workers init

	// init datapools
	dataPools := make(map[string]schema.DataPool)
	for _, sh := range config.Schemas {
		dataPools[sh.Name] = schema.NewDataPool(sh)
	}

	// todo: in a parallel depending on type
	for _, job := range config.Jobs {
		func() {
			// todo: fix here, no schema data pool will be nill
			dataPool := dataPools[job.Schema]
			worker, error := NewWorker(config, job, dataPool)
			if error != nil {
				panic("Worker initialization error")
			}
			defer worker.Close()
			worker.InitIntervalReportingSummary()
			worker.Work()
			worker.Summary()
			worker.ExtendCopySavedFieldsToDataPool()
		}()
	}
}

type Worker struct {
	cfg         *config.Config
	job         *config.Job
	wg          sync.WaitGroup
	db          database.Client
	handler     JobHandler
	rateLimiter Limiter
	pool        JobPool
	dataPool    schema.DataPool
	Report      Report
	ticker      *time.Ticker
	startTime   time.Time
}

func NewWorker(cfg *config.Config, job *config.Job, dataPool schema.DataPool) (*Worker, error) {
	// todo: check errors
	fmt.Printf("Starting job: %s\n", lo.If(job.Name != "", job.Name).Else(job.Type))
	worker := new(Worker)
	worker.cfg = cfg
	worker.job = job
	worker.wg.Add(int(job.Connections))
	worker.Report = NewReport(job)
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
	w.startTime = time.Now()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	go func() {
		<-interruptChan
		w.Cancel()
	}()

	for i := 0; i < int(w.job.Connections); i++ {
		go func() {
			defer w.wg.Done()
			for w.pool.SpawnJob() {
				w.rateLimiter.Take()
				// perform operation
				duration, error := w.handler.Handle()
				w.Report.Add(duration, error)
				// add debug of some kind
				if error != nil {
					// todo: debug
					log.Debug(error)
				}
				w.pool.MarkJobDone()
			}
		}()
	}
	w.wg.Wait()
	w.Report.SetDuration(time.Since(w.startTime))
}

func (w *Worker) InitIntervalReportingSummary() {
	reportingFormat := w.job.GetReport()
	if reportingFormat == nil || reportingFormat.Interval == 0 {
		log.Info("Interval reporting skipped")
		return
	}

	w.ticker = time.NewTicker(reportingFormat.Interval)
	go func(worker *Worker) {
		for range w.ticker.C {
			worker.Report.SetDuration(time.Since(w.startTime))
			worker.Report.Summary()
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
	w.Report.Summary()
}

func (w *Worker) Cancel() {
	w.pool.Cancel()
	w.Close()
}

func (w *Worker) Close() {
	if w.job.Type != string(config.Sleep) {
		if err := w.db.Disconnect(); err != nil {
			panic(err)
		}
	}
	if w.ticker != nil {
		w.ticker.Stop()
	}
}
