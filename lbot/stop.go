package lbot

import (
	"context"

	"github.com/kuzxnia/loadbot/lbot/proto"
)

type StoppingProcess struct {
	proto.UnimplementedStopProcessServer
	ctx  context.Context
	lbot *Lbot
}

func NewStoppingProcess(ctx context.Context, lbot *Lbot) *StoppingProcess {
	return &StoppingProcess{ctx: ctx, lbot: lbot}
}

func (c *StoppingProcess) Run(ctx context.Context, request *proto.StopRequest) (*proto.StopResponse, error) {
	// validate is configured

	go c.lbot.Cancel()
	// if watch arg - run watch
	return &proto.StopResponse{}, nil
}
