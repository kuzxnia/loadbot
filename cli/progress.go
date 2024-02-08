package cli

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/gosuri/uilive"
	"github.com/kuzxnia/loadbot/lbot/proto"
	"google.golang.org/grpc"
)

func WorkloadProgress(conn grpc.ClientConnInterface, request *proto.ProgressRequest) (err error) {
	client := proto.NewProgressProcessClient(conn)

	stream, err := client.Run(context.TODO(), request)
	if err != nil {
		return fmt.Errorf("starting stress test failed: %w", err)
	}

	writer := uilive.New()
	writer.Start()
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("cannot receive %v", err)
		}
    fmt.Fprintf(writer, "RPS: %.1f Requests: %d ErrorRate %.2f\n", resp.Rps, resp.RequestsTotal, resp.ErrorRate)
	}
	writer.Stop()

	return
}
