package lbot

import (
	"context"
	"fmt"

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
	done := make(chan bool)

	go func() {
		// for message := range w.lbot.logs {
		// 	resp := proto.WatchResponse{Message: message}

		// 	if err := srv.Send(&resp); err != nil {
		// 		// todo: handle client not connected
		// 		log.Printf("client closed connection, closing channel done")
		// 		done <- true
		// 		return
		// 	}
		// }
	}()
	<-done
	fmt.Printf("done")

	return nil
}
