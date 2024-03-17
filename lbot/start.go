package lbot

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kuzxnia/loadbot/lbot/proto"
	"github.com/kuzxnia/loadbot/lbot/worker"
	"github.com/samber/lo"
)

type StartProcess struct {
	proto.UnimplementedStartProcessServer
	ctx  context.Context
	lbot *Lbot
}

func NewStartProcess(ctx context.Context, lbot *Lbot) *StartProcess {
	return &StartProcess{ctx: ctx, lbot: lbot}
}

func (c *StartProcess) Run(ctx context.Context, request *proto.StartRequest) (*proto.StartResponse, error) {
	err := c.lbot.Run()

	return &proto.StartResponse{}, err
}

func (c *StartProcess) RunWithProgress(request *proto.StartWithProgressRequest, srv proto.StartProcess_RunWithProgressServer) error {
	interval, err := time.ParseDuration(request.RefreshInterval)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		c.lbot.Run()
		wg.Done()
	}()

	ticker := time.NewTicker(interval)
	go func() {
		defer wg.Done()
		var notDoneWorkers []*worker.Worker

		for range ticker.C {
			select {
			case <-c.lbot.done:
				fmt.Println("workload done")
				return
			default:
			}
			if notDoneWorkers == nil || len(notDoneWorkers) == 0 {
				notDoneWorkers = lo.Filter(lo.Values(c.lbot.workers), func(worker *worker.Worker, index int) bool {
					return !worker.IsDone()
				})
			}
			for _, w := range notDoneWorkers {
				isWorkerFinished := w.IsDone()
				resp := proto.ProgressResponse{
					Requests:          w.Metrics.Requests(),
					Duration:          uint64(w.Metrics.DurationSeconds()),
					Rps:               w.Metrics.Rps(),
					ErrorRate:         w.Metrics.ErrorRate(),
					IsFinished:        isWorkerFinished,
					JobName:           w.JobName(),
					RequestOperations: w.RequestedOperations(),
					RequestDuration:   w.RequestedDurationSeconds(),
				}
				if err := srv.Send(&resp); err != nil {
					// todo: handle client not connected
					log.Printf("Client closed connection")
					return
				}
				if isWorkerFinished {
					notDoneWorkers = lo.Filter(lo.Values(c.lbot.workers), func(worker *worker.Worker, index int) bool {
						return !worker.IsDone()
					})
				}
			}
		}
	}()

	wg.Wait()

	return nil
}
