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

	fmt.Println("Initializing workers")
	for _, job := range config.Jobs {
		worker, error := NewWorker(config, &job)
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
		go worker.Work()
	}
	fmt.Println("Workers started")
	ml.start = time.Now() // add progress bar if running with limit

	ml.wg.Wait()
	ml.Summary()
}

func (ml *mongoload) Summary() {
	for _, worker := range ml.workers {
		worker.Summary()
	}
}

func (ml *mongoload) cancel() {
	for _, worker := range ml.workers {
		worker.Cancel()
	}
}

type Worker struct {
	wg          sync.WaitGroup
	db          database.Client
	handler     JobHandler
	rateLimiter Limiter
	pool        JobPool
	Statistic   Stats
	startTime   time.Time
}

func NewWorker(cfg *config.Config, job *config.Job) (*Worker, error) {
	db, err := database.NewMongoClient(cfg.ConnectionString, job, &cfg.Schemas[0])
	if err != nil {
		return nil, err
	}

	// todo: check errors
	worker := new(Worker)
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
}

func (w *Worker) ExecuteJob() {
	defer w.wg.Done()

	for w.pool.SpawnJob() {
    w.rateLimiter.Take()
		// perform operation
		start := time.Now()
		// do sth with is error
		_, error := w.handler.Handle()
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

func (w *Worker) Summary() {
	elapsed := time.Since(w.startTime)
	requestsDone := w.pool.GetRequestsDone()
	rps := float64(requestsDone) / elapsed.Seconds()

	fmt.Printf("\nTime took %f s\n", elapsed.Seconds())
	fmt.Printf("Total operations: %d\n", requestsDone)
	fmt.Printf("Requests per second: %f rp/s\n", rps)
	w.Statistic.Summary()

	// if batch := w.db.GetBatchSize(); batch > 1 {
	// 	fmt.Printf("Operations per second: %f op/s\n", float64(requestsDone*batch)/elapsed.Seconds())
	// }
}

func (w *Worker) Cancel() {
	w.pool.Cancel()
}
