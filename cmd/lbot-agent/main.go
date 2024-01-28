package main

import (
	"context"

	"github.com/kuzxnia/loadbot/lbot"
	"github.com/kuzxnia/loadbot/lbot/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := log.NewLogger(ctx)
	if err != nil {
		panic(err)
	}
	agent := lbot.NewAgent(ctx, logger)

	agent.Listen()
}

// func main() {
// 	// maxprocs.Set()
// 	config, err := args.Parse()
// 	if err != nil {
// 		panic(err)
// 	}
// 	log := logger.Default()
// 	log.SetConfig(config)
// 	defer log.CloseOutputFile()

// 	driver.Torment(config)
// }
