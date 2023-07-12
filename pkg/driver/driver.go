package driver

import (
	"fmt"
	"sync"
	"time"

	"github.com/kuzxnia/mongoload/pkg/config"
	"github.com/kuzxnia/mongoload/pkg/database"
	"github.com/kuzxnia/mongoload/pkg/logger"
)

var log = logger.Default()

type mongoload struct {
	config  *config.Config
	db      database.Client
	wg      sync.WaitGroup
	workers []Worker
	start   time.Time
}

// todo: change params to options struct
// todo: move database part to worker 
func New(config *config.Config, database database.Client) (*mongoload, error) {
	load := new(mongoload)
	load.config = config
	load.db = database

	fmt.Println("Initializing workers")
	for _, job := range config.Jobs {
		worker, error := NewWorker(&job)
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

// func (ml *mongoload) performWriteOperation() (bool, error) {
// 	start := time.Now()
// 	writedSuccessfuly, error := ml.db.InsertOneOrMany()
// 	elapsed := time.Since(start)
// 	ml.writeStats.Add(float64(elapsed.Milliseconds()), error)

// 	// add debug of some kind
// 	if error != nil {
// 		// todo: debug
// 		log.Debug(error)
// 	}

// 	return writedSuccessfuly, error
// }

// func (ml *mongoload) performReadOperation() (bool, error) {
// 	start := time.Now()
// 	writedSuccessfuly, error := ml.db.ReadOne()
// 	elapsed := time.Since(start)
// 	ml.readStats.Add(float64(elapsed.Milliseconds()), error)
// 	if error != nil {
// 		// todo: debug
// 		log.Debug(error)
// 	}

// 	return writedSuccessfuly, error
// }

// func (ml *mongoload) performUpdateOperation() (bool, error) {
// 	start := time.Now()
// 	writedSuccessfuly, error := ml.db.UpdateOne()
// 	elapsed := time.Since(start)
// 	ml.updateStats.Add(float64(elapsed.Milliseconds()), error)
// 	if error != nil {
// 		// todo: debug
// 		log.Debug(error)
// 	}

// 	return writedSuccessfuly, error
// }
