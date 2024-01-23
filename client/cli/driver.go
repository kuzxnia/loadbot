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

func provideStartDriverHandler() *cobra.Command {
	cmd := cobra.Command{
		Use:     CommandStartDriver,
		Aliases: []string{"i"},
		Short:   "Start stress test",
		Args:    cobra.ExactArgs(installationArgsNum),
		RunE:    installationHandler,
	}

	// flags := cmd.Flags()
	// flags

	return &cmd
}

func startDriverHandler(cmd *cobra.Command, args []string) error {
	// flags := cmd.Flags()

	request := driver.StartRequest{}

	logger.Info("🚀 Starting stress test")

	if err := driver.NewStartProcess(cmd.Context(), request).Run(); err != nil {
		return fmt.Errorf("starting stress test failed: %w", err)
	}

	logger.Info("✅ Starting stress test succeeded")

	return nil
}

func provideStopDriverHandler() *cobra.Command {
	cmd := cobra.Command{
		Use:     CommandStopDriver,
		Aliases: []string{"i"},
		Short:   "Stop stress test",
		Args:    cobra.ExactArgs(installationArgsNum),
		RunE:    installationHandler,
	}

	return &cmd
}

func stopDriverHandler(cmd *cobra.Command, args []string) error {
	request := driver.StopRequest{}

	logger.Info("🚀 Stop stress test")

	if err := driver.NewStopProcess(cmd.Context(), request).Run(); err != nil {
		return fmt.Errorf("stop stress test failed: %w", err)
	}

	logger.Info("✅ Stop stress test succeeded")

	return nil
}

func provideWatchDriverHandler() *cobra.Command {
	cmd := cobra.Command{
		Use:     CommandWatchDriver,
		Aliases: []string{"i"},
		Short:   "Watch stress test",
		Args:    cobra.ExactArgs(installationArgsNum),
		RunE:    installationHandler,
	}

	return &cmd
}

func watchDriverHandler(cmd *cobra.Command, args []string) error {
	request := driver.WatchRequest{}

	logger.Info("🚀 Watch stress test")

	if err := driver.NewWatchProcess(cmd.Context(), request).Run(); err != nil {
		return fmt.Errorf("Watch stress test failed: %w", err)
	}

	logger.Info("✅ Watch stress test succeeded")

	return nil
}
