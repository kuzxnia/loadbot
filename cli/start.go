package cli

func startingDriverHandler(cmd *cobra.Command, args []string) error {
	// flags := cmd.Flags()

	request := StartRequest{}

	logger.Info("🚀 Starting stress test")

	if err := NewStartProcess(cmd.Context(), request).Run(); err != nil {
		return fmt.Errorf("starting stress test failed: %w", err)
	}

	logger.Info("✅ Starting stress test succeeded")

	return nil
}
