package cli

import (
	"context"

	"github.com/kuzxnia/loadbot/lbot/proto"
	"google.golang.org/grpc"
)

func StopDriver(conn grpc.ClientConnInterface, request *proto.StopRequest) (err error) {
	Logger.Info("🚀 Stopping stress test")

	client := proto.NewStopProcessClient(conn)

	reply, err := client.Run(context.TODO(), request)
	if err != nil {
		Logger.Fatal("arith error:", err)
		return
	}

	Logger.Infof("Received: %v", reply)
	Logger.Info("✅ Stopping stress test succeeded")

	return nil
}
