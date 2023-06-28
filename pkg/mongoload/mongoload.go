package mongoload

import (
	"fmt"
	"sync"
	"time"

	"github.com/kuzxnia/mongoload/pkg/database"
	"github.com/kuzxnia/mongoload/pkg/rps"
)

type Book struct {
	Title  string
	Author string
	ISBN   string
}

type mongoload struct {
	db database.DbClient
	wg sync.WaitGroup

	concurrentConnections int
	rateLimit             int // rps limit
	operationsAmount      int64

	// rlp      sync.Mutex
	// requests int64

	start       time.Time
	duration    time.Duration
	rateLimiter rps.Limiter

	jobs chan int
	done chan bool
}

func New(ops int, conns int, rateLimit int, duration time.Duration, database database.DbClient) (*mongoload, error) {
	load := new(mongoload)
	// uri, rps, time, ops, conns,
	if duration != 0 && ops != 0 {
		panic("duration or ops are required")
	} else if duration != 0 {
		load.duration = duration
		load.operationsAmount = 0
	} else {
		load.operationsAmount = int64(ops)
		load.duration = 0
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

	load.wg.Add(load.concurrentConnections)
	load.jobs = make(chan int)
	load.done = make(chan bool, load.operationsAmount)

	return load, nil
}

func (ml *mongoload) Torment() {
	// start workers
	for i := 0; i < ml.concurrentConnections; i++ {
		go ml.worker()
	}
	fmt.Println("workers started")

	ml.start = time.Now()
	// send jobs
	for i := 0; i < int(ml.operationsAmount); i++ {
		ml.jobs <- i
	}

	close(ml.jobs)
	ml.wg.Wait()
	close(ml.done)

	// for result := range ml.done {
	//   fmt.Printf("%v \n", result)
	// }

	// start := time.Now()
	elapsed := time.Since(ml.start)
	fmt.Printf("Time took %f s\n", elapsed.Seconds())
	fmt.Printf("Total operations: %d\n", ml.operationsAmount)
	opsPerSecond := float64(ml.operationsAmount) / elapsed.Seconds()
	fmt.Printf("Ops per second: %f ops/s\n", opsPerSecond)
}

func (ml *mongoload) worker() {
	// fmt.Println("worker starting")
	defer ml.wg.Done()

	for range ml.jobs {
    ml.rateLimiter.Take()
		ml.done <- ml.performSingleWrite()
	}
}

func (ml *mongoload) performSingleWrite() bool {
	writedSuccessfuly, _ := ml.db.InsertOne()
	// handle error in stats -> change '_' from above
	return writedSuccessfuly
}
func (*mongoload) performSingleRead() {}
