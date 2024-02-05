package cli

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/kuzxnia/loadbot/lbot/proto"
	"google.golang.org/grpc"
)

func WatchDriver(conn grpc.ClientConnInterface, request *proto.WatchRequest) (err error) {
	Logger.Info("🚀 Starting stress test")

	client := proto.NewWatchProcessClient(conn)

	stream, err := client.Run(context.TODO(), request)
	if err != nil {
		return fmt.Errorf("starting stress test failed: %w", err)
	}

	done := make(chan bool)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true // means stream is finished
				return
			}
			if err != nil {
				log.Fatalf("cannot receive %v", err)
			}
			log.Printf("%s", resp.Message)
		}
	}()

	<-done // we will wait until all response is received

	Logger.Info("✅ Starting stress test succeeded")

	return
}
