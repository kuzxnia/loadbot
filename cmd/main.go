package main

import (
	"github.com/kuzxnia/mongoload/pkg/args"
	"github.com/kuzxnia/mongoload/pkg/driver"
	"github.com/kuzxnia/mongoload/pkg/logger"
)

func main() {
	// maxprocs.Set()
	config, err := args.Parse()
	if err != nil {
		panic(err)
	}
	log := logger.Default()
	log.SetConfig(config)
	defer log.CloseOutputFile()

	driver, error := driver.New(config)
	if error != nil {
		panic(error)
	}
	driver.Torment()
}
