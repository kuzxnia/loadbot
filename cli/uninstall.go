package cli

import (
	"net/rpc"

	"github.com/kuzxnia/loadbot/orchiestrator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func unInstallationHandler(cmd *cobra.Command, args []string) error {
	request := orchiestrator.UnInstallationRequest{}

	log.Info("🚀 Starting installation process")

	var reply int
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("Found errors trying to connect to lbot-agent:", err)
	}

	err = client.Call("UnInstallationProcess.Run", request, &reply)
	if err != nil {
		log.Fatal("UnInstallationProcess error:", err)
		// 	return fmt.Errorf("installation failed: %w", err)
	}

	log.Info("✅ Installation process succeeded")

	return nil
}
