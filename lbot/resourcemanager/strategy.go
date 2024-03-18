package resourcemanager

import (
	"time"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
)

const (
	LocalDockerStrategy = "docker"
	HelmChartStrategy   = "helm"
)

type ResourceManagerConfig struct {
	KubeconfigPath string
	Context        string
	Namespace      string
	HelmTimeout    time.Duration
}

type InstallRequest struct {
	ResourceManagerConfig
	Name                 string
	HelmValues           []string
	WorkloadConfigString string
}

type InstallResponse struct{}

type UpgradeRequest struct {
	ResourceManagerConfig
	Name                 string
	HelmValues           []string
	WorkloadConfigString string
}

type UpgradeResponse struct{}

type UnInstallRequest struct {
	ResourceManagerConfig
	Name string
}

type UnInstallResponse struct{}

type ListRequest struct {
	ResourceManagerConfig
}

type ListResponse struct{}

type ResourceManager interface {
	Install(*InstallRequest) error
	Upgrade(*UpgradeRequest) error
	UnInstall(*UnInstallRequest) error
	List(*ListRequest) error
}

var (
	Strategies = []string{LocalDockerStrategy, HelmChartStrategy}

	nameToStrategy = map[string]ResourceManager{
		// LocalDockerStrategy: &DockerService{},
		HelmChartStrategy: &HelmManager{},
	}

	DefaultStrategy = LocalDockerStrategy
	HelmProviders   = getter.All(cli.New())
)

func GetResourceManager(cfg *ResourceManagerConfig) (ResourceManager, error) {
	manager, err := NewHelmManager(cfg)
	if err != nil {
		// todo: message failed do create manager
		return nil, err
	}
	return ResourceManager(manager), nil

	// return nameToStrategy[DefaultStrategy], nil
}
