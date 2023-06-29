package main

import (
	"flag"

	"github.com/kuzxnia/mongoload/pkg/database"
	"github.com/kuzxnia/mongoload/pkg/mongoload"
)

func main() {
	// iserts
	mongoURI := flag.String("uri", "mongodb://localhost:27017", "Database hostname url")
	mongoDatabase := flag.String("db", "load_test", "Database name")
	mongoCollection := flag.String("col", "load_test_coll", "Collection name")
	concurrentConnections := flag.Int("conn", 100, "Concurrent connections amount")
	rpsLimit := flag.Int("rps", 0, "RPS limit")
	durationLimit := flag.Duration("d", 0, "Duration limit")
	opsAmount := flag.Int("req", 0, "Requests to perform")
	batchSize := flag.Uint64("bs", 0, "Batch size")
	dataLenght := flag.Uint64("dl", 100, "Lenght of single item data(chars)")

	flag.Parse()

	poolSize := *concurrentConnections * 8
	db, err := database.NewMongoDbClient(*mongoURI, *mongoDatabase, *mongoCollection, uint64(poolSize), *batchSize, *dataLenght)

	defer func() {
		if err = db.Disconnect(); err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	load, _ := mongoload.New(
		*opsAmount,
		*concurrentConnections,
		*rpsLimit,
		*durationLimit,
		db,
	)

	load.Torment()
}
