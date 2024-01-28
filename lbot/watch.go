package lbot

import "context"

type WatchingRequest struct{}

type WatchingProcess struct {
	ctx context.Context
}

func NewWatchingProcess(ctx context.Context) *WatchingProcess {
	return &WatchingProcess{ctx: ctx}
}

func (w *WatchingProcess) Run(request *WatchingRequest, reply *int) error {
	// if watch arg - run watch
	return nil
}
