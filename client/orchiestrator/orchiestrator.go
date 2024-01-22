package orchiestrator

type Orchiestrator struct {
	chart *HelmChart
}

func NewOrchiestrator() (*Orchiestrator, error) {
	chart, err := NewHelmChart()
	if err != nil {
		return nil, err
	}

	return &Orchiestrator{
		chart: chart,
	}, nil
}

func (o *Orchiestrator) createResources() {
  o.chart.install()
}

func (o *Orchiestrator) deleteResources() {
  o.chart.uninstall()
}
