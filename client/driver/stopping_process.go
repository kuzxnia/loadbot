package driver

import (
	"context"
)

type StoppingRequest struct{}

type StoppingProcess struct {
	ctx     context.Context
	request *StoppingRequest
}

func NewStoppingProcess(ctx context.Context, request StoppingRequest) *StoppingProcess {
	return &StoppingProcess{ctx: ctx, request: &request}
}

func (c *StoppingProcess) Run() error {
	// if watch arg - run watch
	return nil
}
