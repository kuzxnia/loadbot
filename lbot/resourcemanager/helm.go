package resourcemanager

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/kuzxnia/loadbot/lbot/k8s"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli/values"
)

//go:embed workload-chart.tgz
var chartBytes []byte

// defult are from above,  - MVP
// but it should be able to process helm charts from internet also

type HelmManager struct {
	cfg           *ResourceManagerConfig
	chart         *chart.Chart
	clusterClient *k8s.ClusterClient
}

// add optional argument with chart version
func NewHelmManager(cfg *ResourceManagerConfig) (*HelmManager, error) {
	// use default or fetch from internet from tag
	// todo: later add validation for type

	chart, err := loader.LoadArchive(bytes.NewReader(chartBytes))
	if err != nil {
		return nil, err
	}

	clusterClient, err := k8s.GetClusterClient(cfg.KubeconfigPath, cfg.Context)
	if err != nil {
		return nil, err
	}

	return &HelmManager{
		cfg:           cfg,
		chart:         chart,
		clusterClient: clusterClient,
	}, nil
}

func (c *HelmManager) Install(request *InstallRequest) (err error) {
	// 1. write values to file

	// 2. helm action config
	// namespace, release, timout, kube config, context
	installConfig := new(action.Configuration)
	installConfig.Init(
		c.clusterClient.RESTClientGetter,
		c.cfg.Namespace,
		os.Getenv("HELM_DRIVER"),
		log.Printf,
	)

	// 3. installer
	installer := action.NewInstall(installConfig)
	installer.Namespace = request.Namespace
	installer.ReleaseName = "dummy-release-name"
	installer.Timeout = c.cfg.HelmTimeout
	installer.Labels["role"] = "workload"

	// 4. get cli values
	options := values.Options{
		Values: []string{"workload.name=" + request.Name},
	}

	vals, err := options.MergeValues(HelmProviders)
	if err != nil {
		return err
	}
	// 5. merge them with helm value file

	// 5. install
	if _, err = installer.Run(c.chart, vals); err != nil {
		return fmt.Errorf("failed to install helm chart: %w", err)
	}

	return
}

func (c *HelmManager) UnInstall(request *UnInstallRequest) (err error) {
	cfg := new(action.Configuration)
	cfg.Init(
		c.clusterClient.RESTClientGetter,
		c.cfg.Namespace,
		os.Getenv("HELM_DRIVER"),
		log.Printf,
	)
	uninstaller := action.NewUninstall(cfg)

	uninstaller.Run(request.Name)
	return
	// todo: add are you shure you want to uninstall sth?
}

func (c *HelmManager) Upgrade(request *UpgradeRequest) (err error) {
	cfg := new(action.Configuration)
	cfg.Init(
		c.clusterClient.RESTClientGetter,
		c.cfg.Namespace,
		os.Getenv("HELM_DRIVER"),
		log.Printf,
	)
	upgrader := action.NewUpgrade(cfg)
	upgrader.Namespace = request.Namespace
	upgrader.Timeout = c.cfg.HelmTimeout
	upgrader.Labels["role"] = "workload"

	options := values.Options{
		Values: []string{"workload.name=" + request.Name},
	}

	vals, err := options.MergeValues(HelmProviders)
	if err != nil {
		return err
	}

	if _, err = upgrader.Run(request.Name, c.chart, vals); err != nil {
		return fmt.Errorf("failed to install helm chart: %w", err)
	}
  return
}

func (c *HelmManager) List(*ListRequest) (err error) {
	cfg := new(action.Configuration)
	cfg.Init(
		c.clusterClient.RESTClientGetter,
		c.cfg.Namespace,
		os.Getenv("HELM_DRIVER"),
		log.Printf,
	)

	list := action.NewList(cfg)
	list.Selector = "role=workload"

	releases, err := list.Run()
	for _, release := range releases {
		fmt.Println(release.Name, release.Namespace, release.Info.Description)
	}

	return
}
