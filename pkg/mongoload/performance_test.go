package mongoload_test 

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/kuzxnia/mongoload/pkg/config"
	"github.com/kuzxnia/mongoload/pkg/mongoload"
)

type DummyClient struct{}

func (c *DummyClient) Disconnect() error {
	return nil
}

func (c *DummyClient) InsertOne() (bool, error) {
	return true, nil
}

func (c *DummyClient) InsertMany() (bool, error) {
	return true, nil
}

func (c *DummyClient) InsertOneOrMany() (bool, error) {
	// time.Sleep(time.Millisecond * 10)
	return true, nil
}

func (c *DummyClient) ReadOne() (bool, error) {
	return true, nil
}

func (c *DummyClient) ReadMany() (bool, error) {
	return true, nil
}

func (c *DummyClient) GetBatchSize() uint64 {
	return 0
}

func benchmarkMongoload(c *config.Config, b *testing.B) {
	ml, error := mongoload.New(c, &DummyClient{})
	if error != nil {
		b.Error(error)
	}
	fmt.Println(runtime.NumCPU())
	b.SetParallelism(int(c.ConcurrentConnections) / runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ml.Torment()
		}
	})
	ml.Summary()
	// summary
}

func BenchmarkMongoloadMaxRps(b *testing.B) {
	benchmarkMongoload(&config.Config{
		MongoURI:              "",
		MongoDatabase:         "",
		MongoCollection:       "",
		ConcurrentConnections: 100,
		PoolSize:              100,
		RpsLimit:              0,
		DurationLimit:         0,
		OpsAmount:             10000,
		BatchSize:             0,
		DataLenght:            1000,
		WriteRatio:            0.5,
		Timeout:               time.Second,
	}, b)
}
