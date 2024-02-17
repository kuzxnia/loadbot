package main

import (
	"context"
	"os"

	"github.com/kuzxnia/loadbot/cli"
	log "github.com/sirupsen/logrus"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	rootCmd := cli.New(version, commit, date)

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		log.Errorf("❌ Error: %s", err.Error())
		return 1
	}

	return 0
}
