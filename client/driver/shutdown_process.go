package driver

import (
	"context"
)

type StopRequest struct{}

type StopProcess struct {
	ctx     context.Context
	request *StopRequest
}

func NewStopProcess(ctx context.Context, request StopRequest) *StopProcess {
	return &StopProcess{ctx: ctx, request: &request}
}

func (c *StopProcess) Run() error {
	// if watch arg - run watch
	return nil
}
