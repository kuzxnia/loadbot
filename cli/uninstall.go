package cli

import (
	"net/rpc"

	"github.com/kuzxnia/loadbot/orchiestrator"
	"github.com/spf13/cobra"
)

func unInstallationHandler(cmd *cobra.Command, args []string) error {
	request := orchiestrator.UnInstallationRequest{}

	Logger.Info("🚀 Starting installation process")

	var reply int
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		Logger.Fatal("Found errors trying to connect to lbot-agent:", err)
	}

	err = client.Call("UnInstallationProcess.Run", request, &reply)
	if err != nil {
		Logger.Fatal("UnInstallationProcess error:", err)
		// 	return fmt.Errorf("installation failed: %w", err)
	}

	Logger.Info("✅ Installation process succeeded")

	return nil
}
