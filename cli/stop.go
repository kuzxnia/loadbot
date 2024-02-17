package cli

import (
	"context"

	"github.com/kuzxnia/loadbot/lbot/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func StopDriver(conn grpc.ClientConnInterface, request *proto.StopRequest) (err error) {
	log.Info("🚀 Stopping stress test")

	client := proto.NewStopProcessClient(conn)

	reply, err := client.Run(context.TODO(), request)
	if err != nil {
		log.Fatal("arith error:", err)
		return
	}

	log.Infof("Received: %v", reply)
	log.Info("✅ Stopping stress test succeeded")

	return nil
}
