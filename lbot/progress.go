package lbot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kuzxnia/loadbot/lbot/proto"
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
	done := make(chan bool)

	interval, err := time.ParseDuration(request.RefreshInterval)
	if err != nil {
		return err
	}

	// stats shouls be assigned to job/ worker
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
      // add job name fileter
			for _, worker := range w.lbot.workers {
				resp := proto.ProgressResponse{
          // add job name here
					Requests:   worker.Metrics.Requests(),
					Rps:        worker.Metrics.Rps(),
					ErrorRate:  worker.Metrics.ErrorRate(),
					Operations: worker.RequestedOperations(),
					Duration:   worker.RequestedDurationSeconds(),
				}

				if err := srv.Send(&resp); err != nil {
					// todo: handle client not connected
					log.Printf("client closed connection, closing channel done")
					done <- true
					return
				}
			}
		}
	}()

	<-done
	fmt.Printf("done")

	return nil
}
