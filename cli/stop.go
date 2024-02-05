package cli

import (
	"context"

	"github.com/kuzxnia/loadbot/lbot/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func stoppingDriverHandler(cmd *cobra.Command, args []string) (err error) {
	Logger.Info("🚀 Stopping stress test")

	conn, err := grpc.Dial("127.0.0.1:1235", grpc.WithInsecure())
	if err != nil {
		Logger.Fatal("Found errors trying to connect to lbot-agent:", err)
		return
	}
	defer conn.Close()
	client := proto.NewStopProcessClient(conn)

	request := proto.StopRequest{}

	reply, err := client.Run(context.TODO(), &request)
	if err != nil {
		Logger.Fatal("arith error:", err)
		return
	}

	Logger.Infof("Received: %v", reply)
	Logger.Info("✅ Stopping stress test succeeded")

	return nil
}
