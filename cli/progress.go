package cli

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/kuzxnia/loadbot/lbot/proto"
	"google.golang.org/grpc"
)

func WorkloadProgress(conn grpc.ClientConnInterface, request *proto.ProgressRequest) (err error) {
	client := proto.NewProgressProcessClient(conn)

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
			log.Printf("rsp=%d total=%d err_rate=%.2f", resp.Rps, resp.RequestsTotal, resp.ErrorRate)
		}
	}()

	<-done // we will wait until all response is received

	return
}
