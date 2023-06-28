package mongoload

import (
	"fmt"
	"sync"
	"time"

	"github.com/kuzxnia/mongoload/pkg/database"
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

	// rps
	// rlp      sync.Mutex
	// requests int64
	// start time.Time

	jobs chan int
	done chan bool
}

func New(ops int, conns int, rps int, time int, database database.DbClient) (*mongoload, error) {
	load := new(mongoload)
	// uri, rps, time, ops, conns,
	if (time == 0 && rps == 0) && ops == 0 {
		panic("time and rps or ops are required")
	} else if time != 0 && rps != 0 {
		load.operationsAmount = int64(time * rps)
	} else {
		load.operationsAmount = int64(ops)
	}
	load.rateLimit = rps

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
	start := time.Now()
	for i := 0; i < ml.concurrentConnections; i++ {
		go ml.worker()
	}
	fmt.Println("workers started")

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
	elapsed := time.Since(start)
	fmt.Printf("Time took %f s\n", elapsed.Seconds())
	fmt.Printf("Total operations: %d\n", ml.operationsAmount)
	opsPerSecond := float64(ml.operationsAmount) / elapsed.Seconds()
	fmt.Printf("Ops per second: %f ops/s\n", opsPerSecond)
}

func (ml *mongoload) worker() {
	// fmt.Println("worker starting")
	defer ml.wg.Done()

	for range ml.jobs {
		ml.done <- ml.performSingleWrite()
	}
}

func (ml *mongoload) performSingleWrite() bool {
	writedSuccessfuly, _ := ml.db.InsertOne()
	// handle error in stats -> change '_' from above
	return writedSuccessfuly
}
func (*mongoload) performSingleRead() {}
