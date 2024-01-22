package orchiestrator

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

type HelmChart struct {
	chart *chart.Chart
}

// add optional argument with chart version
func NewHelmChart() (*HelmChart, error) {
	// use default or fetch from internet from tag
	// todo: later add validation for type

	chart, err := loader.LoadArchive(bytes.NewReader(chartBytes))
	if err != nil {
		return nil, err
	}

	return &HelmChart{
		chart: chart,
	}, nil
}

func (c *HelmChart) install() (err error) {
	return
}

func (c *HelmChart) uninstall() (err error) {
  return
}
