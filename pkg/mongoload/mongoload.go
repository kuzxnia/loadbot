package mongoload

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
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

	readHistogram  Histogram
	writeHistogram Histogram
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
	load.readHistogram = NewHistogram()
	load.writeHistogram = NewHistogram()

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

	// wp := ml.stats.WriteHistogram.Percentiles([]float64{0.5, 0.9, 0.99})
	// rp := ml.stats.ReadHistogram.Percentiles([]float64{0.5, 0.9, 0.99})

	fmt.Printf("\nTime took %f s\n", elapsed.Seconds())
	fmt.Printf("Total operations: %d\n", requestsDone)

	if ml.writeRatio != 0 {
		wmean, _ := ml.writeHistogram.Mean()
		wp50, _ := ml.writeHistogram.Percentile(50)
		wp90, _ := ml.writeHistogram.Percentile(90)
		wp99, _ := ml.writeHistogram.Percentile(99)
		fmt.Printf("Total writes: %d\n", ml.writeHistogram.Len())
		fmt.Printf("Write AVG: %.2fms, P50: %.2fms, P90: %.2fms P99: %f.2ms\n", wmean, wp50, wp90, wp99)
	}

	if ml.writeRatio != 1 {
		rmean, _ := ml.readHistogram.Mean()
		rp50, _ := ml.readHistogram.Percentile(50)
		rp90, _ := ml.readHistogram.Percentile(90)
		rp99, _ := ml.readHistogram.Percentile(99)
		fmt.Printf("Total reads: %d\n", ml.readHistogram.Len())
		fmt.Printf("Read AVG: %.2fms, P50: %.2fms, P90: %.2fms P99: %.2fms\n", rmean, rp50, rp90, rp99)
	}

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
			ml.performWriteOperation()
		} else {
			ml.performReadOperation()
		}
		ml.pool.MarkJobDone()
	}
}

func (ml *mongoload) performWriteOperation() bool {
	start := time.Now()
	writedSuccessfuly, _ := ml.db.InsertOneOrMany()
	elapsed := time.Since(start)
	ml.writeHistogram.Update(float64(elapsed.Milliseconds()))

	// handle error in stats -> change '_' from above
	return writedSuccessfuly
}

func (ml *mongoload) performReadOperation() bool {
	start := time.Now()
	writedSuccessfuly, _ := ml.db.ReadOne()
	elapsed := time.Since(start)
	ml.readHistogram.Update(float64(elapsed.Milliseconds()))

	// handle error in stats -> change '_' from above
	return writedSuccessfuly
}
