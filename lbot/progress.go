package lbot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kuzxnia/loadbot/lbot/driver"
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

	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			resp := proto.ProgressResponse{
				RequestsTotal: driver.Stats.Total(),
				Rps:           driver.Stats.Rps(),
				ErrorRate:     driver.Stats.ErrorRate(),
			}

			if err := srv.Send(&resp); err != nil {
				// todo: handle client not connected
				log.Printf("client closed connection, closing channel done")
				done <- true
				return
			}
		}
	}()

	<-done
	fmt.Printf("done")

	return nil
}
