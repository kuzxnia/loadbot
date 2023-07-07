package main

import (
	"github.com/kuzxnia/mongoload/cmd/args"
	"github.com/kuzxnia/mongoload/pkg/database"
	"github.com/kuzxnia/mongoload/pkg/driver"
	"github.com/kuzxnia/mongoload/pkg/logger"
)

func main() {
	config, err := args.Parse()
	if err != nil {
		panic(err)
	}
	log := logger.Default()
	log.SetConfig(config)
	defer log.CloseOutputFile()

	db, err := database.NewMongoClient(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = db.Disconnect(); err != nil {
			panic(err)
		}
	}()

	load, error := driver.New(config, db)
	if error != nil {
		panic(error)
	}
	load.Torment()
}
