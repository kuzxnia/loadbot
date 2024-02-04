package cli

import (
	"net/rpc"

	"github.com/kuzxnia/loadbot/lbot"
	"github.com/spf13/cobra"
)

func watchingDriverHandler(cmd *cobra.Command, args []string) error {
	request := lbot.WatchingRequest{}

	Logger.Info("🚀 Watching stress test")

	var reply int
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		Logger.Fatal("Found errors trying to connect to lbot-agent:", err)
	}

  // currently not possible because it require stream
  // for first version change this to execute in look, and fetch last n min/sec of logs
	err = client.Call("WatchProcess.Run", request, &reply)
	if err != nil {
		Logger.Fatal("WatchProcess error:", err)
	}

	Logger.Info("✅ Watching stress test succeeded")

	return nil
}
