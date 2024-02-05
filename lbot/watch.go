package lbot

import (
	"context"
	"fmt"
	"log"
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

  // todo: only temporary flush channel
	drainchan(w.lbot.logs)

	done := make(chan bool)

	go func() {
		for {
			//   select {
			//   case:
			// }
			// read channel with logs
			time.Sleep(time.Second)
			select {
			case message := <-w.lbot.logs:
				resp := proto.WatchResponse{Message: message}

				if err := srv.Send(&resp); err != nil {
					// todo: handle client not connected
					log.Printf("client closed connection, closing channel done")
					done <- true
					return
				}
			}
		}
		// todo: or do this by interatin over channel
	}()

	fmt.Printf("before done")

	<-done
	fmt.Printf("done")

	return nil
}

func drainchan(chann chan string) {
	for {
		select {
		case <-chann:
		default:
			return
		}
	}
}
