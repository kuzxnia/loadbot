package driver

import (
	"context"
)

type StartRequest struct {
	watch bool
}

type StartProcess struct {
	ctx     context.Context
	request *StartRequest
}

func NewStartProcess(ctx context.Context, request StartRequest) *StartProcess {
	return &StartProcess{ctx: ctx, request: &request}
}

func (c *StartProcess) Run() error {
	// if watch arg - run watch
	return nil
}
