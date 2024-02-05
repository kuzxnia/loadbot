package lbot

import (
	"context"

	"github.com/kuzxnia/loadbot/lbot/proto"
)

type StartProcess struct {
	proto.UnimplementedStartProcessServer
	ctx  context.Context
	lbot *Lbot
}

func NewStartProcess(ctx context.Context, lbot *Lbot) *StartProcess {
	return &StartProcess{ctx: ctx, lbot: lbot}
}

func (c *StartProcess) Run(ctx context.Context, request *proto.StartRequest) (*proto.StartResponse, error) {
	// if watch arg - run watch

	// validate is configured
	go c.lbot.Run()

	// before starting process it will varify health of cluster, if pods
	return &proto.StartResponse{}, nil
}
