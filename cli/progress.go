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

type ProgressBar struct {
	bars map[string]*pb.ProgressBar
}

// todo: no need to send request ops and request duration in every request
// pull job data in init
func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		bars: make(map[string]*pb.ProgressBar),
	}
}

func (b *ProgressBar) IsInitialized(resp *proto.ProgressResponse) bool {
	if lo.IsNil(resp) {
		return len(b.bars) > 0
	}
	_, ok := b.bars[resp.JobName]
	return ok
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
	bar.Set("job", resp.JobName)
	bar.Set("requestOperations", resp.RequestOperations)
	bar.Set("requestDuration", resp.RequestDuration)

	b.bars[resp.JobName] = bar
}

func (b *ProgressBar) Start(resp *proto.ProgressResponse) {
	b.bars[resp.JobName].Start()
}

func (b *ProgressBar) Update(resp *proto.ProgressResponse) {
	bar := b.bars[resp.JobName]
	if resp.RequestDuration != 0 {
		bar.SetCurrent(int64(resp.GetDuration()))
	} else if resp.RequestOperations != 0 {
		bar.SetCurrent(int64(resp.GetRequests()))
	}
	bar.Set("rps", int(resp.GetRps()))
	bar.Set("requests", resp.GetRequests())
	bar.Set("duration", resp.GetDuration())

	bar.Write()

	if resp.IsFinished {
		bar.Finish()
		fmt.Println() // flush stdout
	}
}

func (b *ProgressBar) Finish() {
	for _, bar := range b.bars {
		bar.Finish()
	}
}
