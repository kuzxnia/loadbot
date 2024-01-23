package driver

import (
	"context"
)

type WatchRequest struct{}

type WatchProcess struct {
	ctx     context.Context
	request *WatchRequest
}

func NewWatchProcess(ctx context.Context, request WatchRequest) *WatchProcess {
	return &WatchProcess{ctx: ctx, request: &request}
}

func (c *WatchProcess) Run() error {
	// if watch arg - run watch
	return nil
}
