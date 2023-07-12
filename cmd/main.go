package main

import (
	"fmt"

	"github.com/kuzxnia/mongoload/pkg/args"
)

func main() {
	// maxprocs.Set()
	fmt.Println("hello")

	cfg, err := args.Parse()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", cfg)
	// log := logger.Default()
	// log.SetConfig(config)
	// defer log.CloseOutputFile()

	// db, err := database.NewMongoClient(config)
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() {
	// 	if err = db.Disconnect(); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// load, error := driver.New(config, db)
	// if error != nil {
	// 	panic(error)
	// }
	// load.Torment()
}
