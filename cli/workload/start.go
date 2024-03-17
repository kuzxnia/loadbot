package workload 

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/kuzxnia/loadbot/lbot/proto"
	"google.golang.org/grpc"
)

// checks if process is running in local system

// tutaj nie powinno wchodzić proto
func StartWorkload(conn grpc.ClientConnInterface, request *proto.StartRequest) (err error) {
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

func StartWorkloadWithProgress(conn grpc.ClientConnInterface, request *proto.StartWithProgressRequest) (err error) {
	// todo: mapowanie to proto
	fmt.Println("🚀 Starting stress test")

	client := proto.NewStartProcessClient(conn)

	stream, err := client.RunWithProgress(context.TODO(), request)
	if err != nil {
		return fmt.Errorf("starting stress test failed: %w", err)
	}

	fmt.Println("✅ Starting stress test succeeded")

	bar := NewProgressBar()
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("cannot receive %v", err)
		}

		if !bar.IsInitialized(resp) {
			bar.Init(resp)
			bar.Start(resp)
		}

		bar.Update(resp)
	}

	if bar.IsInitialized(nil) {
		bar.Finish()
	} else {
		// in that case no response was received - no job running
		fmt.Println("There are no running jobs")
	}
	return
}
