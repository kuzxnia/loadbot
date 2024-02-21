package cli

import (
	"context"
	"fmt"

	"github.com/kuzxnia/loadbot/lbot/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func StopDriver(conn grpc.ClientConnInterface, request *proto.StopRequest) (err error) {
	fmt.Println("🚀 Stopping stress test")

	client := proto.NewStopProcessClient(conn)

	_, err = client.Run(context.TODO(), request)
	if err != nil {
		log.Fatal("arith error:", err)
		return
	}

	fmt.Println("✅ Stopping stress test succeeded")

	return nil
}
