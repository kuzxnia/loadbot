package main

import (
	"github.com/kuzxnia/loadbot/lbot/pkg/args"
	"github.com/kuzxnia/loadbot/lbot/pkg/driver"
	"github.com/kuzxnia/loadbot/lbot/pkg/logger"
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

	driver.Torment(config)
}
