package cli

import (
	"context"
	"fmt"

	"github.com/kuzxnia/loadbot/lbot/proto"
	"google.golang.org/grpc"
)

// checks if process is running in local system

// tutaj nie powinno wchodzić proto
func StartDriver(conn grpc.ClientConnInterface, request *proto.StartRequest) (err error) {
	// todo: mapowanie to proto
	fmt.Println("🚀 Starting stress test")

	client := proto.NewStartProcessClient(conn)

	_, err = client.Run(context.TODO(), request)
	if err != nil {
		return fmt.Errorf("starting stress test failed: %w", err)
	}

	fmt.Println("✅ Starting stress test succeeded")

	return
}
