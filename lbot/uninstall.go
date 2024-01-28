package cli

import (
	"context"
	"time"
)

type UnInstallationRequest struct {
	KubeconfigPath string
	Context        string
	Namespace      string
	HelmTimeout    time.Duration
}

type UnInstallationProcess struct {
	ctx     context.Context
	request *UnInstallationRequest
}

func NewUnInstallationProcess(ctx context.Context, request UnInstallationRequest) *UnInstallationProcess {
	return &UnInstallationProcess{ctx: ctx, request: &request}
}

func (c *UnInstallationProcess) Run() error {
	// if watch arg - run watch
	return nil
}
