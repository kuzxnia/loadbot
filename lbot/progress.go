package lbot

import (
	"context"
	"fmt"
	"time"

	"github.com/kuzxnia/loadbot/lbot/driver"
	"github.com/kuzxnia/loadbot/lbot/proto"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
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
	// refactor
	if len(lo.Filter(w.lbot.workers, func(worker *driver.Worker, index int) bool { return !worker.IsDone() })) == 0 {
		log.Printf("There are no running jobs")
		return nil
	}

	done := make(chan bool)
	ticker := time.NewTicker(interval)
	go func() {
		notDoneWorkers := lo.Filter(w.lbot.workers, func(worker *driver.Worker, index int) bool {
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
					notDoneWorkers = lo.Filter(w.lbot.workers, func(worker *driver.Worker, index int) bool {
						return !worker.IsDone()
					})
				}
				select {
				case <-w.lbot.done:
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
