package lbot

import (
	"context"
	"log"
	"time"

	"github.com/kuzxnia/loadbot/lbot/driver"
	"github.com/kuzxnia/loadbot/lbot/proto"
	"github.com/samber/lo"
)

type ProgressProcess struct {
	proto.UnimplementedProgressProcessServer
	ctx  context.Context
	lbot *Lbot
}

func NewProgressProcess(ctx context.Context, lbot *Lbot) *ProgressProcess {
	return &ProgressProcess{ctx: ctx, lbot: lbot}
}

func (w *ProgressProcess) Run(request *proto.ProgressRequest, srv proto.ProgressProcess_RunServer) error {
	interval, err := time.ParseDuration(request.RefreshInterval)
	if err != nil {
		return err
	}

	runningWorkers := lo.Filter(w.lbot.workers, func(worker *driver.Worker, index int) bool {
		return !worker.IsDone()
	})

	if len(runningWorkers) == 0 {
		log.Printf("There are no running jobs")
		return nil
	}

	done := make(chan bool)
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			// todo: add job name fileter
			for _, worker := range runningWorkers {
				isWorkerFinished := worker.IsDone()
				resp := proto.ProgressResponse{
					Requests:          worker.Metrics.Requests(),
					Duration:          uint64(worker.Metrics.DurationSeconds()),
					Rps:               worker.Metrics.Rps(),
					ErrorRate:         worker.Metrics.ErrorRate(),
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
					log.Printf("Worker finished running jobs")
					done <- true
					return
				}
			}
		}
	}()
	<-done

	return nil
}
