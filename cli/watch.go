package cli

import "context"

type WatchingRequest struct{}

type WatchingProcess struct {
	ctx     context.Context
	request *WatchingRequest
}

func NewWatchingProcess(ctx context.Context, request WatchingRequest) *WatchingProcess {
	return &WatchingProcess{ctx: ctx, request: &request}
}

func (w *WatchingProcess) Run() error {
	// if watch arg - run watch
	return nil
}
