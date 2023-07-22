package driver

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/kuzxnia/mongoload/pkg/config"
	"github.com/kuzxnia/mongoload/pkg/database"
	"github.com/kuzxnia/mongoload/pkg/logger"
)

var log = logger.Default()

type mongoload struct {
	config  *config.Config
	wg      sync.WaitGroup
	workers []*Worker
	start   time.Time
}

// todo: change params to options struct
// todo: move database part to worker
func New(config *config.Config) (*mongoload, error) {
	load := new(mongoload)
	load.config = config

	// todo: ping db, before workers init

	// todo: now all jobs will be executed in a parallel,
	// change this to be execexuted as queue or in a parallel depending on type
	load.wg.Add(len(config.Jobs))

	fmt.Println("Initializing workers")
	for _, job := range config.Jobs {
		worker, error := NewWorker(config, job)
		if error != nil {
			panic("Worker initialization error")
		}
		load.workers = append(load.workers, worker)
	}
	fmt.Println("Workers initialized")

	return load, nil
}

func (ml *mongoload) Torment() {
	for _, worker := range ml.workers {
		func(worker *Worker) {
			defer ml.wg.Done()
			worker.Work()
      worker.Summary()
		}(worker)
	}

	ml.start = time.Now() // add progress bar if running with limit

	ml.wg.Wait()
}

func (ml *mongoload) Summary() {
	for _, worker := range ml.workers {
		worker.Summary()
	}
}

func (ml *mongoload) Cancel() {
	for _, worker := range ml.workers {
		worker.Cancel()
	}
}

func (ml *mongoload) Close() {
	for _, worker := range ml.workers {
		worker.Close()
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
	Report      Report
}

func NewWorker(cfg *config.Config, job *config.Job) (*Worker, error) {
	// todo: check errors
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

	worker.handler = NewJobHandler(job, worker.db)
	// todo: init db
	return worker, nil
}

func (w *Worker) Work() {
	startTime := time.Now()

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
	w.Report.SetDuration(time.Since(startTime))
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
}
