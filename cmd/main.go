package main

import (
	"fmt"

	"github.com/kuzxnia/mongoload/pkg/args"
	"github.com/kuzxnia/mongoload/pkg/driver"
	"github.com/kuzxnia/mongoload/pkg/logger"
)

func main() {
	// maxprocs.Set()
	fmt.Println("hello")

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
