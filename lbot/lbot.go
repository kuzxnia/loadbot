package lbot

import (
	"context"

	"github.com/kuzxnia/loadbot/lbot/config"
	"github.com/kuzxnia/loadbot/lbot/driver"
	"github.com/kuzxnia/loadbot/lbot/schema"
)

type Lbot struct {
	ctx     context.Context
	config  *config.Config
	workers []*driver.Worker
	logs    chan string
	metric  *driver.Metric
}

func NewLbot(ctx context.Context) *Lbot {
	return &Lbot{
		ctx:    ctx,
		metric: driver.NewMetrics(),
		logs:   make(chan string),
	}
}

func (l *Lbot) Run() {
	// todo: ping db, before workers init

	// init datapools
	dataPools := make(map[string]schema.DataPool)
	for _, sh := range l.config.Schemas {
		dataPools[sh.Name] = schema.NewDataPool(sh)
	}

	// // todo: in a parallel depending on type
	for _, job := range l.config.Jobs {
		func() {
			// todo: fix here, no schema data pool will be nill
			dataPool := dataPools[job.Schema]
			worker, error := driver.NewWorker(l.ctx, l.config, job, dataPool, l.metric)
			l.workers = append(l.workers, worker)
			if error != nil {
				panic("Worker initialization error")
			}
			defer worker.Close()
			worker.InitIntervalReportingSummary(l.logs)
			worker.Work()
			worker.Summary()
			worker.ExtendCopySavedFieldsToDataPool()
		}()
	}
}

func (l *Lbot) Cancel() error {
	for _, worker := range l.workers {
		worker.Cancel()
	}
	l.workers = nil

	return nil
}

func (l *Lbot) SetConfig(config *config.Config) {
	l.config = config
}
