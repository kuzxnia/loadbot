package lbot

import (
	"time"

	"github.com/kuzxnia/loadbot/lbot/proto"
	"github.com/kuzxnia/loadbot/lbot/resourcemanager"
	"golang.org/x/net/context"
)

type Orchiestrator struct {
	proto.UnimplementedOrchistratorServiceServer
	ctx context.Context
}

func NewOrchiestrator(ctx context.Context) *Orchiestrator {
	return &Orchiestrator{ctx: ctx}
}

func (o *Orchiestrator) Install(ctx context.Context, request *proto.InstallRequest) (*proto.InstallResponse, error) {
	// if watch arg - run watch

	timeout, err := time.ParseDuration(request.HelmTimeout)
	if err != nil {
		return nil, err
	}
	cfg := resourcemanager.ResourceManagerConfig{
		KubeconfigPath: request.KubeconfigPath,
		Context:        request.Context,
		Namespace:      request.Namespace,
		HelmTimeout:    timeout,
	}

	resourceManager, err := resourcemanager.GetResourceManager(&cfg)
	if err != nil {
		return nil, err
	}

	err = resourceManager.Install(request)

	// create resources,

	// for helm
	// set config map thought values or helm file values - configure thought config map
	// - change to yaml values, it will be driver settings and workload test setting

	// for docker it will be save inside container??

	// if we have multiple nodes, we need to set cluster from them

	// if flag starting is provided it will start workload
	// same with watch flag

	return &proto.InstallResponse{}, err
}

func (o *Orchiestrator) UnInstall(ctx context.Context, request *proto.UnInstallRequest) (*proto.UnInstallResponse, error) {
	timeout, err := time.ParseDuration(request.HelmTimeout)
	if err != nil {
		return nil, err
	}
	cfg := resourcemanager.ResourceManagerConfig{
		KubeconfigPath: request.KubeconfigPath,
		Context:        request.Context,
		Namespace:      request.Namespace,
		HelmTimeout:    timeout,
	}

	resourceManager, err := resourcemanager.GetResourceManager(&cfg)
	if err != nil {
		return nil, err
	}

	err = resourceManager.UnInstall(request)

	return &proto.UnInstallResponse{}, err
}
