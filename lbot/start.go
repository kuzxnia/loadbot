package lbot

import "golang.org/x/net/context"

type StartRequest struct {
	Watch bool
}

type StartProcess struct {
	ctx  context.Context
	lbot *Lbot
}

func NewStartProcess(ctx context.Context, lbot *Lbot) *StartProcess {
	return &StartProcess{ctx: ctx, lbot: lbot}
}

func (c *StartProcess) Run(request *StartRequest, reply *int) error {
	// if watch arg - run watch

  // validate is configured

	go c.lbot.Run()

	// before starting process it will varify health of cluster, if pods
	return nil
}
