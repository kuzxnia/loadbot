package mongoload

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/kuzxnia/mongoload/pkg/config"
	"github.com/kuzxnia/mongoload/pkg/database"
	"github.com/kuzxnia/mongoload/pkg/logger"
	"github.com/kuzxnia/mongoload/pkg/rps"
	"github.com/kuzxnia/mongoload/pkg/worker"
)

var log = logger.Default()

type mongoload struct {
	config *config.Config
	db     database.Client
	wg     sync.WaitGroup

	start       time.Time
	duration    time.Duration
	rateLimiter rps.Limiter

	pool worker.JobPool

	readStats   Stats
	writeStats  Stats
	updateStats Stats
}

// todo: change params to options struct
func New(config *config.Config, database database.Client) (*mongoload, error) {
	load := new(mongoload)

	load.config = config

	if config.DurationLimit == 0 && config.OpsAmount == 0 {
		load.pool = worker.NewNoLimitTimerJobPool()
	} else if config.DurationLimit != 0 {
		load.duration = config.DurationLimit
		load.pool = worker.NewTimerJobPool(config.DurationLimit)
	} else {
		load.pool = worker.NewDeductionJobPool(uint64(load.config.OpsAmount))
	}

	// change to is pointer nil
	if config.RpsLimit == 0 {
		load.rateLimiter = rps.NewNoLimitLimiter()
	} else {
		load.rateLimiter = rps.NewSimpleLimiter(config.RpsLimit)
	}

	load.db = database
	load.readStats = NewReadStats()
	load.writeStats = NewWriteStats()
	load.updateStats = NewUpdateStats()

	load.wg.Add(int(load.config.ConcurrentConnections))

	return load, nil
}

func (ml *mongoload) Torment() {
	// handle interrupt
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	go func() {
		<-interruptChan
		ml.cancel()
	}()

	fmt.Println("Starting workers")
	for i := 0; i < int(ml.config.ConcurrentConnections); i++ {
		go ml.worker()
	}
	fmt.Println("Workers started")
	ml.start = time.Now()
	// add progress bar if running with limit

	ml.wg.Wait()
	ml.Summary()
}

func (ml *mongoload) Summary() {
	elapsed := time.Since(ml.start)

	requestsDone := ml.pool.GetRequestsDone()
	rps := float64(requestsDone) / elapsed.Seconds()

	fmt.Printf("\nTime took %f s\n", elapsed.Seconds())
	fmt.Printf("Total operations: %d\n", requestsDone)

	ml.writeStats.Summary()
	ml.readStats.Summary()
	ml.updateStats.Summary()

	fmt.Printf("Requests per second: %f rp/s\n", rps)
	if batch := ml.db.GetBatchSize(); batch > 1 {
		fmt.Printf("Operations per second: %f op/s\n", float64(requestsDone*batch)/elapsed.Seconds())
	}
}

func (ml *mongoload) cancel() {
	print("\nCancelling...")
	ml.pool.Cancel()
}

func (ml *mongoload) worker() {
	defer ml.wg.Done()

	for ml.pool.SpawnJob() {

		ml.rateLimiter.Take()
		requestTypeFactor := ml.pool.GetRequestsStarted() % 100

		if requestTypeFactor < ml.config.WriteRatio {
			ml.performWriteOperation()
		} else if requestTypeFactor < ml.config.ReadRatio {
			ml.performReadOperation()
		} else if requestTypeFactor < ml.config.UpdateRatio {
			ml.performUpdateOperation()
		}

		ml.pool.MarkJobDone()
	}
}

func (ml *mongoload) performWriteOperation() (bool, error) {
	start := time.Now()
	writedSuccessfuly, error := ml.db.InsertOneOrMany()
	elapsed := time.Since(start)
	ml.writeStats.Add(float64(elapsed.Milliseconds()), error)

	// add debug of some kind
	if error != nil {
		// todo: debug
		log.Debug(error)
	}

	return writedSuccessfuly, error
}

func (ml *mongoload) performReadOperation() (bool, error) {
	start := time.Now()
	writedSuccessfuly, error := ml.db.ReadOne()
	elapsed := time.Since(start)
	ml.readStats.Add(float64(elapsed.Milliseconds()), error)
	if error != nil {
		// todo: debug
		log.Debug(error)
	}

	return writedSuccessfuly, error
}

func (ml *mongoload) performUpdateOperation() (bool, error) {
	start := time.Now()
	writedSuccessfuly, error := ml.db.UpdateOne()
	elapsed := time.Since(start)
	ml.updateStats.Add(float64(elapsed.Milliseconds()), error)
	if error != nil {
		// todo: debug
		log.Debug(error)
	}

	return writedSuccessfuly, error
}
