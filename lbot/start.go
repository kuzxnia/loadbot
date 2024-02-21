package lbot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kuzxnia/loadbot/lbot/driver"
	"github.com/kuzxnia/loadbot/lbot/proto"
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
	// if watch arg - run watch

	// validate is configured
	go c.lbot.Run()

	// before starting process it will varify health of cluster, if pods
	return &proto.StartResponse{}, nil
}

func (c *StartProcess) RunWithProgress(request *proto.StartWithProgressRequest, srv proto.StartProcess_RunWithProgressServer) error {
	// if watch arg - run watch

	// validate is configured
	go c.lbot.Run()

	interval, err := time.ParseDuration(request.RefreshInterval)
	if err != nil {
		return err
	}
	// Setup a 1-minute timeout for the job to start
	timeout := time.After(1 * time.Minute)
	jobStarted := make(chan bool)

	go func() {
		for {
			if len(lo.Filter(c.lbot.workers, func(worker *driver.Worker, index int) bool { return !worker.IsDone() })) > 0 {
				jobStarted <- true
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	select {
	case <-jobStarted:
		// Job started, continue with the rest of the function
	case <-timeout:
		log.Printf("Timeout occurred: No job started within 1 minute")
	}

	done := make(chan bool)
	ticker := time.NewTicker(interval)
	go func() {
		notDoneWorkers := lo.Filter(c.lbot.workers, func(worker *driver.Worker, index int) bool {
			return !worker.IsDone()
		})
		for range ticker.C {
			for _, worker := range notDoneWorkers {
				isWorkerFinished := worker.IsDone()
				resp := proto.ProgressResponse{
					Requests:          worker.Metrics.Requests(),
					Duration:          uint64(worker.Metrics.DurationSeconds()),
					Rps:               worker.Metrics.Rps(),
					ErrorRate:         worker.Metrics.ErrorRate(),
					IsFinished:        isWorkerFinished,
					JobName:           worker.JobName(),
					RequestOperations: worker.RequestedOperations(),
					RequestDuration:   worker.RequestedDurationSeconds(),
				}
				if err := srv.Send(&resp); err != nil {
					// todo: handle client not connected
					log.Printf("Client closed connection")
					done <- true
					return
				}
				if isWorkerFinished {
					notDoneWorkers = lo.Filter(c.lbot.workers, func(worker *driver.Worker, index int) bool {
						return !worker.IsDone()
					})
				}
				select {
				case <-c.lbot.done:
					fmt.Println("workload done")
					done <- true
				default:
				}

			}
		}
	}()
	<-done

	return nil
}
