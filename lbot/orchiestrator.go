package lbot

import (
	"github.com/kuzxnia/loadbot/lbot/resourcemanager"
	"golang.org/x/net/context"
)

type Orchiestrator struct {
	ctx context.Context
}

func NewOrchiestrator(ctx context.Context) *Orchiestrator {
	return &Orchiestrator{ctx: ctx}
}

func (o *Orchiestrator) Install(ctx context.Context, request *resourcemanager.InstallRequest) (*resourcemanager.InstallResponse, error) {
	cfg := resourcemanager.ResourceManagerConfig{
		KubeconfigPath: request.KubeconfigPath,
		Context:        request.Context,
		Namespace:      request.Namespace,
		HelmTimeout:    request.HelmTimeout,
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

	return &resourcemanager.InstallResponse{}, err
}

func (o *Orchiestrator) UnInstall(ctx context.Context, request *resourcemanager.UnInstallRequest) (*resourcemanager.UnInstallResponse, error) {
	cfg := resourcemanager.ResourceManagerConfig{
		KubeconfigPath: request.KubeconfigPath,
		Context:        request.Context,
		Namespace:      request.Namespace,
		HelmTimeout:    request.HelmTimeout,
	}

	resourceManager, err := resourcemanager.GetResourceManager(&cfg)
	if err != nil {
		return nil, err
	}

	err = resourceManager.UnInstall(request)

	return &resourcemanager.UnInstallResponse{}, err
}

func (o *Orchiestrator) Upgrade(ctx context.Context, request *resourcemanager.UpgradeRequest) (*resourcemanager.UpgradeResponse, error) {
	cfg := resourcemanager.ResourceManagerConfig{
		KubeconfigPath: request.KubeconfigPath,
		Context:        request.Context,
		Namespace:      request.Namespace,
		HelmTimeout:    request.HelmTimeout,
	}

	resourceManager, err := resourcemanager.GetResourceManager(&cfg)
	if err != nil {
		return nil, err
	}

	err = resourceManager.Upgrade(request)

	return &resourcemanager.UpgradeResponse{}, err
}

func (o *Orchiestrator) List(ctx context.Context, request *resourcemanager.ListRequest) (*resourcemanager.ListResponse, error) {
	cfg := resourcemanager.ResourceManagerConfig{
		KubeconfigPath: request.KubeconfigPath,
		Context:        request.Context,
		Namespace:      request.Namespace,
		HelmTimeout:    request.HelmTimeout,
	}

	resourceManager, err := resourcemanager.GetResourceManager(&cfg)
	if err != nil {
		return nil, err
	}

	err = resourceManager.List(request)

	return &resourcemanager.ListResponse{}, nil
}
