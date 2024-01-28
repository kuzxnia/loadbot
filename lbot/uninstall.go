package lbot

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
	ctx context.Context
}

func NewUnInstallationProcess(ctx context.Context) *UnInstallationProcess {
	return &UnInstallationProcess{ctx: ctx}
}

func (c *UnInstallationProcess) Run(request *UnInstallationRequest, reply *int) error {
	// if watch arg - run watch
	return nil
}
