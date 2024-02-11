package resourcemanager

import (
	"bytes"
	_ "embed"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

//go:embed helm-chart.tgz
var chartBytes []byte

// defult are from above,  - MVP
// but it should be able to process helm charts from internet also

type HelmService struct {
	chart *chart.Chart
}

// add optional argument with chart version
func NewHelmChart() (*HelmService, error) {
	// use default or fetch from internet from tag
	// todo: later add validation for type

	chart, err := loader.LoadArchive(bytes.NewReader(chartBytes))
	if err != nil {
		return nil, err
	}

	return &HelmService{
		chart: chart,
	}, nil
}

func (c *HelmService) Install() (err error) {
	return
}

func (c *HelmService) UnInstall() (err error) {
	return
}

func (c *HelmService) Suspend() (err error) {
	return
}
