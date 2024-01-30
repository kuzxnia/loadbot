package lbot

import (
	"context"

	"github.com/kuzxnia/loadbot/lbot/config"
	"github.com/kuzxnia/loadbot/lbot/driver"
	"github.com/kuzxnia/loadbot/lbot/schema"
)

type Lbot struct {
	config *config.Config
}

func NewLbot(config *config.Config) *Lbot {
	return &Lbot{
		config: config,
	}
}

func (l *Lbot) Run(ctx context.Context) {
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
			worker, error := driver.NewWorker(l.config, job, dataPool)
			if error != nil {
				panic("Worker initialization error")
			}
			defer worker.Close()
			worker.InitIntervalReportingSummary()
			worker.Work()
			worker.Summary()
			worker.ExtendCopySavedFieldsToDataPool()
		}()
	}
}

func (l *Lbot) Ping() error {
	return nil
}

func (l *Lbot) SetConfig(config *ConfigRequest) {
	// todo: create config from request
	// should be in config
	// l.config = config
}
