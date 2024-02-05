package cli

import (
	"context"
	"fmt"

	"github.com/kuzxnia/loadbot/lbot/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// checks if process is running in local system

func startingDriverHandler(cmd *cobra.Command, args []string) (err error) {
	Logger.Info("🚀 Starting stress test")

	conn, err := grpc.Dial("127.0.0.1:1235", grpc.WithInsecure())
	if err != nil {
		Logger.Fatal("Found errors trying to connect to lbot-agent:", err)
		return
	}
	defer conn.Close()

	client := proto.NewStartProcessClient(conn)

	request := proto.StartRequest{
		Watch: false,
	}

	reply, err := client.Run(context.TODO(), &request)

	if err != nil {
		return fmt.Errorf("starting stress test failed: %w", err)
	}

  Logger.Infof("Received: %v", reply)
	Logger.Info("✅ Starting stress test succeeded")

	return
}
