package lbot

import "context"

type WatchingRequest struct{}

type WatchingProcess struct {
	ctx  context.Context
	lbot *Lbot
}

func NewWatchingProcess(ctx context.Context, lbot *Lbot) *WatchingProcess {
	return &WatchingProcess{ctx: ctx, lbot: lbot}
}

func (w *WatchingProcess) Run(request *WatchingRequest, reply *int) error {
	// if watch arg - run watch
	return nil
}
