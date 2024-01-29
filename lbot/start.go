package lbot

import "golang.org/x/net/context"

type StartRequest struct {
	watch bool
}

type StartProcess struct {
	ctx context.Context
}

func NewStartProcess(ctx context.Context) *StartProcess {
	return &StartProcess{
		ctx: ctx,
	}
}

func (c *StartProcess) Run(request *StartRequest, reply *int) error {
	// if watch arg - run watch

  // 	driver.Torment(config)

	// before starting process it will varify health of cluster, if pods
	return nil
}
