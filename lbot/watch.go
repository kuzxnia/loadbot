package lbot

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kuzxnia/loadbot/lbot/proto"
)

type WatchingRequest struct{}

type WatchingProcess struct {
	proto.UnimplementedWatchProcessServer
	ctx  context.Context
	lbot *Lbot
}

func NewWatchingProcess(ctx context.Context, lbot *Lbot) *WatchingProcess {
	return &WatchingProcess{ctx: ctx, lbot: lbot}
}

func (w *WatchingProcess) Run(request *proto.WatchRequest, srv proto.WatchProcess_RunServer) error {
	// if watch arg - run watch
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(count int64) {
			defer wg.Done()

			// time sleep to simulate server process time
			time.Sleep(time.Duration(count) * time.Second)
			resp := proto.WatchResponse{Message: fmt.Sprintf("Request #%d", count)}
			if err := srv.Send(&resp); err != nil {
				log.Printf("send error %v", err)
        // todo: handle client not connected 
			}
			log.Printf("finishing request number : %d", count)
		}(int64(i))
	}

	wg.Wait()
	return nil
}
