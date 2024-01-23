package orchiestrator

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
	request *UnInstallationProcess
}

func NewUnInstallationProcess(ctx context.Context, request UnInstallationProcess) *UnInstallationProcess {
	return &UnInstallationProcess{ctx: ctx, request: &request}
}

func (c *UnInstallationProcess) Run() error {
	// if watch arg - run watch
	return nil
}