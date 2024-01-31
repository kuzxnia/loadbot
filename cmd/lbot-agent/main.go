package main

import (
	"context"
	"os"

	"github.com/kuzxnia/loadbot/lbot"
	"github.com/kuzxnia/loadbot/lbot/log"
)

var (
	// will be overridden by goreleaser: https://goreleaser.com/cookbooks/using-main.version/
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if exitCode := run(); exitCode != 0 {
		os.Exit(exitCode)
	}
}

func run() int {
	// maxprocs.Set()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := log.NewLogger(ctx)
	if err != nil {
		panic(err)
	}

	cmd := lbot.BuildArgs(logger, version, commit, date)
	err = cmd.ExecuteContext(ctx)
	if err != nil {
		logger.Errorf("❌ Error: %s", err.Error())
		return 1
	}
	return 0
}
