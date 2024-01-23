package installation

import (
	"context"
	"time"
)

type Request struct {
	KubeconfigPath   string
	Context          string
	Namespace        string
	HelmTimeout      time.Duration
	HelmValuesFiles  []string
	HelmValues       []string
	HelmFileValues   []string
	HelmStringValues []string
}

type InstallationProcess struct {
	ctx     context.Context
	request *Request
}

func NewInstallationProcess(ctx context.Context, request Request) *InstallationProcess {
	return &InstallationProcess{ctx: ctx, request: &request}
}

func (c *InstallationProcess) Run() error {
	// if watch arg - run watch
	return nil
}
