package main

import (
	"github.com/kuzxnia/mongoload/cmd/args"
	"github.com/kuzxnia/mongoload/pkg/database"
	"github.com/kuzxnia/mongoload/pkg/logger"
	"github.com/kuzxnia/mongoload/pkg/mongoload"
)

func main() {
	config, err := args.Parse()
	if err != nil {
		panic(err)
	}
	logger.Default().SetConfig(config)

	db, err := database.NewMongoClient(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = db.Disconnect(); err != nil {
			panic(err)
		}
	}()

	load, error := mongoload.New(config, db)
	if error != nil {
		panic(error)
	}
	load.Torment()
}
