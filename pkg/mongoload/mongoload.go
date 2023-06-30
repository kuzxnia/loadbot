package mongoload

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kuzxnia/mongoload/pkg/database"
	"github.com/kuzxnia/mongoload/pkg/rps"
	"github.com/kuzxnia/mongoload/pkg/worker"
)

type mongoload struct {
	db database.DbClient
	wg sync.WaitGroup

	concurrentConnections int
	rateLimit             int // rps limit
	operationsAmount      int64
	writeRatio            float64

	start       time.Time
	duration    time.Duration
	rateLimiter rps.Limiter

	pool worker.JobPool

	readRequest   uint64
	insertRequest uint64
}

// todo: change params to options struct
func New(
	ops int,
	conns int,
	rateLimit int,
	duration time.Duration,
	database database.DbClient,
	writeRatio float64,
) (*mongoload, error) {
	load := new(mongoload)

	if duration == 0 && ops == 0 {
		load.pool = worker.NewNoLimitTimerJobPool()
	} else if duration != 0 {
		load.duration = duration
		load.pool = worker.NewTimerJobPool(duration)
	} else {
		load.operationsAmount = int64(ops)
		load.pool = worker.NewDeductionJobPool(uint64(load.operationsAmount))
	}

	load.rateLimit = rateLimit
	if rateLimit == 0 {
		load.rateLimiter = rps.NewNoLimitLimiter()
	} else {
		load.rateLimiter = rps.NewSimpleLimiter(rateLimit)
	}

	if conns == 0 {
		conns = 100
	}
	load.concurrentConnections = conns
	load.db = database
	load.writeRatio = writeRatio
	load.readRequest = 0
	load.insertRequest = 0

	load.wg.Add(load.concurrentConnections)

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
	for i := 0; i < ml.concurrentConnections; i++ {
		go ml.worker()
	}
	fmt.Println("Workers started")
	ml.start = time.Now()

	ml.wg.Wait()

	elapsed := time.Since(ml.start)

	requestsDone := ml.pool.GetRequestsDone()
	rps := float64(requestsDone) / elapsed.Seconds()

	fmt.Printf("\nTime took %f s\n", elapsed.Seconds())
	fmt.Printf("Total operations: %d\n", requestsDone)
	fmt.Printf("Total writes: %d\n", atomic.LoadUint64(&ml.insertRequest))
	fmt.Printf("Total reads: %d\n", atomic.LoadUint64(&ml.readRequest))
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
		GetRequestsStarted := ml.pool.GetRequestsStarted()

		ml.rateLimiter.Take()
		if int(GetRequestsStarted)%10 < int(ml.writeRatio*10) {
			ml.performReadOperation()
		} else {
			ml.performInsertOperation()
		}
		ml.pool.MarkJobDone()
	}
}

func (ml *mongoload) performInsertOperation() bool {
	writedSuccessfuly, _ := ml.db.InsertOneOrMany()

	atomic.AddUint64(&ml.insertRequest, 1)
	// if writedSuccessfuly {
	//   fmt.Printf("s")
	// } else {
	//   fmt.Printf("f")
	// }

	// handle error in stats -> change '_' from above
	return writedSuccessfuly
}

func (ml *mongoload) performReadOperation() bool {
	writedSuccessfuly, _ := ml.db.ReadOne()

	atomic.AddUint64(&ml.readRequest, 1)
	// if writedSuccessfuly {
	//   fmt.Printf("s")
	// } else {
	//   fmt.Printf("f")
	// }

	// handle error in stats -> change '_' from above
	return writedSuccessfuly
}
