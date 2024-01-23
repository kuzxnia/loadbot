package cli

import (
	"fmt"

	"github.com/kuzxnia/loadcli/driver"
	"github.com/spf13/cobra"
)

const (
	CommandStartDriver = "start"
	CommandStopDriver  = "stop"
	CommandWatchDriver = "watch"
)

func provideStartingDriverHandler() *cobra.Command {
	cmd := cobra.Command{
		Use:     CommandStartDriver,
		Aliases: []string{"i"},
		Short:   "Start stress test",
		Args:    cobra.ExactArgs(installationArgsNum),
		RunE:    startingDriverHandler,
	}

	// flags := cmd.Flags()
	// flags

	return &cmd
}

func startingDriverHandler(cmd *cobra.Command, args []string) error {
	// flags := cmd.Flags()

	request := driver.StartRequest{}

	logger.Info("🚀 Starting stress test")

	if err := driver.NewStartProcess(cmd.Context(), request).Run(); err != nil {
		return fmt.Errorf("starting stress test failed: %w", err)
	}

	logger.Info("✅ Starting stress test succeeded")

	return nil
}

func provideStoppingDriverHandler() *cobra.Command {
	cmd := cobra.Command{
		Use:     CommandStopDriver,
		Aliases: []string{"i"},
		Short:   "Stopping stress test",
		Args:    cobra.ExactArgs(installationArgsNum),
		RunE:    stoppingDriverHandler,
	}

	return &cmd
}

func stoppingDriverHandler(cmd *cobra.Command, args []string) error {
	request := driver.StoppingRequest{}

	logger.Info("🚀 Stopping stress test")

	if err := driver.NewStoppingProcess(cmd.Context(), request).Run(); err != nil {
		return fmt.Errorf("Stopping stress test failed: %w", err)
	}

	logger.Info("✅ Stopping stress test succeeded")

	return nil
}

func provideWatchingDriverHandler() *cobra.Command {
	cmd := cobra.Command{
		Use:     CommandWatchDriver,
		Aliases: []string{"i"},
		Short:   "Watch stress test",
		Args:    cobra.ExactArgs(installationArgsNum),
		RunE:    watchingDriverHandler,
	}

	return &cmd
}

func watchingDriverHandler(cmd *cobra.Command, args []string) error {
	request := driver.WatchingRequest{}

	logger.Info("🚀 Watching stress test")

	if err := driver.NewWatchingProcess(cmd.Context(), request).Run(); err != nil {
		return fmt.Errorf("Watching stress test failed: %w", err)
	}

	logger.Info("✅ Watching stress test succeeded")

	return nil
}
