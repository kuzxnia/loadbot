package cli

import (
	"fmt"
	"net/rpc"

	"github.com/kuzxnia/loadbot/lbot"
	"github.com/spf13/cobra"
)

// checks if process is running in local system

func startingDriverHandler(cmd *cobra.Command, args []string) (err error) {
	// flags := cmd.Flags()

	request := lbot.StartRequest{
    Watch: false,
  }
	// to change
	var reply int

	Logger.Info("🚀 Starting stress test")

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		Logger.Fatal("Found errors trying to connect to lbot-agent:", err)
		return
	}

	err = client.Call("StartProcess.Run", request, &reply)
	if err != nil {
		return fmt.Errorf("starting stress test failed: %w", err)
	}

	Logger.Info("✅ Starting stress test succeeded")

	return
}
