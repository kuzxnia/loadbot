package lbot

import (
	"context"
)

type StoppingRequest struct{}

type StoppingProcess struct {
	ctx  context.Context
	lbot *Lbot
}

func NewStoppingProcess(ctx context.Context, lbot *Lbot) *StoppingProcess {
	return &StoppingProcess{ctx: ctx, lbot: lbot}
}

func (c *StoppingProcess) Run(request *StoppingRequest, reply *int) error {
	// if watch arg - run watch
	return nil
}
