package lbot

import (
	"context"
	"fmt"
	"time"

	"github.com/kuzxnia/loadbot/lbot/proto"
	"github.com/kuzxnia/loadbot/lbot/worker"
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

func (p *ProgressProcess) Run(request *proto.ProgressRequest, srv proto.ProgressProcess_RunServer) error {
	interval, err := time.ParseDuration(request.RefreshInterval)
	if err != nil {
		return err
	}
	// refactor
	if len(lo.Filter(p.lbot.workers, func(worker *worker.Worker, index int) bool { return !worker.IsDone() })) == 0 {
		log.Printf("There are no running jobs")
		return nil
	}

	done := make(chan bool)
	ticker := time.NewTicker(interval)
	go func() {
		notDoneWorkers := lo.Filter(p.lbot.workers, func(worker *worker.Worker, index int) bool {
			return !worker.IsDone()
		})
		for range ticker.C {
			select {
			case <-p.lbot.done:
				fmt.Println("workload done")
				done <- true
			default:
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
					done <- true
					return
				}
				if isWorkerFinished {
					notDoneWorkers = lo.Filter(p.lbot.workers, func(worker *worker.Worker, index int) bool {
						return !worker.IsDone()
					})
				}

			}
		}
	}()
	<-done

	return nil
}
