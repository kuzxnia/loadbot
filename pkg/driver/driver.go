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
		go func(worker *Worker) {
			defer ml.wg.Done()
			worker.Work()
		}(worker)
	}

	fmt.Println("Workers executed")
	ml.start = time.Now() // add progress bar if running with limit

	ml.wg.Wait()
	ml.Summary()
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
	Statistic   Stats
	startTime   time.Time
}

func NewWorker(cfg *config.Config, job *config.Job) (*Worker, error) {
	db, err := database.NewMongoClient(cfg.ConnectionString, job, job.GetTemplateSchema())
	if err != nil {
		return nil, err
	}

	// todo: check errors
	worker := new(Worker)
	worker.cfg = cfg
	worker.job = job
	worker.wg.Add(int(job.Connections))
	worker.Statistic = NewStatistics(job)
	worker.pool = NewJobPool(job)
	worker.rateLimiter = NewLimiter(job)
	worker.startTime = time.Now()
	worker.db = db
	worker.handler = NewJobHandler(job, db)
	// todo: init db
	return worker, nil
}

func (w *Worker) Work() {
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
				w.Statistic.Add(duration, error)
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
}

func (w *Worker) Summary() {
	elapsed := time.Since(w.startTime)
	requestsDone := w.pool.GetRequestsDone()
	rps := float64(requestsDone) / elapsed.Seconds()

  fmt.Printf("\nJob: \"%s\" took %f s\n", w.job.Name, elapsed.Seconds())
	fmt.Printf("Total operations: %d\n", requestsDone)
	fmt.Printf("Requests per second: %f rp/s\n", rps)
	w.Statistic.Summary()

	// if batch := w.db.GetBatchSize(); batch > 1 {
	// 	fmt.Printf("Operations per second: %f op/s\n", float64(requestsDone*batch)/elapsed.Seconds())
	// }
}

func (w *Worker) Cancel() {
	w.pool.Cancel()
	w.Close()
}

func (w *Worker) Close() {
	if err := w.db.Disconnect(); err != nil {
		panic(err)
	}
}
