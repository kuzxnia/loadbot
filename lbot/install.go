package lbot

import (
	"context"
	"time"
)

type InstallationRequest struct {
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
	ctx context.Context
}

func NewInstallationProcess(ctx context.Context) *InstallationProcess {
	return &InstallationProcess{ctx: ctx}
}

func (c *InstallationProcess) Run(request *InstallationRequest, reply *int) error {
	// if watch arg - run watch

	// create resources,

	// for helm
	// set config map thought values or helm file values - configure thought config map
	// - change to yaml values, it will be driver settings and workload test setting

	// for docker it will be save inside container??

	// if we have multiple nodes, we need to set cluster from them

	// if flag starting is provided it will start workload
	// same with watch flag

	return nil
}
