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
	installConfig := new(action.Configuration)
	installConfig.Init(
		c.clusterClient.RESTClientGetter,
		c.cfg.Namespace,
		os.Getenv("HELM_DRIVER"),
		log.Printf,
	)

	installer := action.NewInstall(installConfig)
	installer.Namespace = request.Namespace
	installer.ReleaseName = request.Name
	installer.Timeout = c.cfg.HelmTimeout
	installer.Labels = map[string]string{"role": "workload"}

	options := values.Options{
		Values:        append([]string{"workload.name=" + request.Name, "workload.namespace=" + request.Namespace}, request.HelmValues...),
		LiteralValues: []string{"workload.config=" + request.WorkloadConfigString},
	}

	vals, err := options.MergeValues(HelmProviders)
	if err != nil {
		return err
	}

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
	// todo: add are you sure you want to uninstall sth?
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
	upgrader.Labels = map[string]string{"role": "workload"}

	options := values.Options{
		Values:        append([]string{"workload.name=" + request.Name, "workload.namespace=" + request.Namespace}, request.HelmValues...),
		LiteralValues: []string{"workload.config=" + request.WorkloadConfigString},
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
