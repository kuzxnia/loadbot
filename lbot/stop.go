package lbot

import (
	"context"
)

type StoppingRequest struct{}

type StoppingProcess struct {
	ctx     context.Context
}

func NewStoppingProcess(ctx context.Context) *StoppingProcess {
	return &StoppingProcess{ctx: ctx}
}

func (c *StoppingProcess) Run(request *StoppingRequest, reply *int) error {
	// if watch arg - run watch
	return nil
}
