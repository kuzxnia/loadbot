package cli

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/kuzxnia/loadbot/lbot/proto"
	"github.com/samber/lo"
	"google.golang.org/grpc"
)

func WorkloadProgress(conn grpc.ClientConnInterface, request *proto.ProgressRequest) (err error) {
	client := proto.NewProgressProcessClient(conn)

	stream, err := client.Run(context.TODO(), request)
	if err != nil {
		return fmt.Errorf("starting stress test failed: %w", err)
	}

	bar := ProgressBar{}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("cannot receive %v", err)
		}

		if !bar.IsInitialized() {
			bar.Init(resp)
      bar.Start()
		}

    bar.Update(resp)
	}


	if bar.IsInitialized() {
		bar.Finish()
	} else {
		// in that case no response was received - no job running
		fmt.Println("There are no running jobs")
	}

	return
}

type ProgressBar struct {
	bar *pb.ProgressBar
}

// todo: no need to send request ops and request duration in every request
// pull job data in init
func NewProgressBar(resp *proto.ProgressResponse) *ProgressBar {
	var value int64
	tmpl := `Job "{{ string . "job" }}" {{ bar . "|" "█" "█" " " "|"}} `
	if resp.GetRequestOperations() != 0 {
		value = int64(resp.GetRequestOperations())
		tmpl += `{{ string . "requests"}}/{{ string . "requestOperations" }}REQ {{string . "rps" }}RPS {{string . "duration"}}S`
	} else {
		value = int64(resp.GetRequestDuration())
		tmpl += `{{ string . "duration"}}/{{ string . "requestDuration" }}S {{string . "rps" }}RPS {{string . "requests"}}REQ`
	}

	bar := pb.New64(int64(value))
	bar.SetTemplateString(tmpl)
	bar.SetWriter(os.Stdout)
	bar.Set(pb.Static, true) // disable auto refresh

	return &ProgressBar{
		bar: bar,
	}
}

func (b *ProgressBar) IsInitialized() bool {
	return !lo.IsNil(b.bar)
}

func (b *ProgressBar) Init(resp *proto.ProgressResponse) {
	var value int64
	tmpl := `Job "{{ string . "job" }}" {{ bar . "|" "█" "█" " " "|"}} `
	if resp.GetRequestOperations() != 0 {
		value = int64(resp.GetRequestOperations())
		tmpl += `{{ string . "requests"}}/{{ string . "requestOperations" }}REQ {{string . "rps" }}RPS {{string . "duration"}}S`
	} else {
		value = int64(resp.GetRequestDuration())
		tmpl += `{{ string . "duration"}}/{{ string . "requestDuration" }}S {{string . "rps" }}RPS {{string . "requests"}}REQ`
	}

	bar := pb.New64(int64(value))
	bar.SetTemplateString(tmpl)
	bar.SetWriter(os.Stdout)
	bar.Set(pb.Static, true) // disable auto refresh
	bar.Set("job", "Insert test")
	bar.Set("requestOperations", resp.RequestOperations)
	bar.Set("requestDuration", resp.RequestDuration)

	b.bar = bar
}

func (b *ProgressBar) Start() {
	b.bar.Start()
}

func (b *ProgressBar) Update(resp *proto.ProgressResponse) {
	if resp.RequestDuration != 0 {
		b.bar.SetCurrent(int64(resp.GetDuration()))
	} else if resp.RequestOperations != 0 {
		b.bar.SetCurrent(int64(resp.GetRequests()))
	}
	b.bar.Set("rps", int(resp.GetRps()))
	b.bar.Set("requests", resp.GetRequests())
	b.bar.Set("duration", resp.GetDuration())

	b.bar.Write()
}

func (b *ProgressBar) Finish() {
  b.bar.Finish()
}
