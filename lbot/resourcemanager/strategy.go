package resourcemanager

import (
	"time"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"

	"github.com/kuzxnia/loadbot/lbot/proto"
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

type ResourceManager interface {
	Install(*proto.InstallRequest) error
	UnInstall(*proto.UnInstallRequest) error
}

var (
	Strategies = []string{LocalDockerStrategy, HelmChartStrategy}

	nameToStrategy = map[string]ResourceManager{
		LocalDockerStrategy: &DockerService{},
		HelmChartStrategy:   &HelmManager{},
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
