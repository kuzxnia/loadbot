package cli

import (
	"net/rpc"

	"github.com/kuzxnia/loadbot/lbot"
	"github.com/spf13/cobra"
)

func stoppingDriverHandler(cmd *cobra.Command, args []string) (err error) {
	request := lbot.StoppingRequest{}

	Logger.Info("🚀 Stopping stress test")

	var reply int

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		Logger.Fatal("Found errors trying to connect to lbot-agent:", err)
		return
	}

	err = client.Call("StoppingProcess.Run", request, &reply)
	if err != nil {
		Logger.Fatal("arith error:", err)
		return
	}

	Logger.Info("✅ Stopping stress test succeeded")

	return nil
}
