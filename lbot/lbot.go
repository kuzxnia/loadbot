package lbot

import (
	"context"
	"fmt"

	"github.com/kuzxnia/loadbot/lbot/config"
	"github.com/kuzxnia/loadbot/lbot/driver"
	"github.com/kuzxnia/loadbot/lbot/schema"
)

type Lbot struct {
	ctx     context.Context
	config  *config.Config
	workers []*driver.Worker
	done    chan bool
}

func NewLbot(ctx context.Context) *Lbot {
	return &Lbot{
		ctx: ctx,
	}
}

func (l *Lbot) Run() {
  l.done = make(chan bool)
	// todo: ping db, before workers init
	// init datapools
	dataPools := make(map[string]schema.DataPool)
	for _, sh := range l.config.Schemas {
		dataPools[sh.Name] = schema.NewDataPool(sh)
	}

	// // todo: in a parallel depending on type
	for _, job := range l.config.Jobs {
		func() {
			dataPool := dataPools[job.Schema]
			worker, error := driver.NewWorker(l.ctx, l.config, job, dataPool)
			if error != nil {
				panic("Worker initialization error")
			}
			fmt.Printf("init worker with job %s\n", job.Name)
			l.workers = append(l.workers, worker)
			// todo: fix here, no schema data pool will be nill
			defer worker.Close()
			worker.InitMetrics()
			worker.Work()
			// worker.Summary()
			worker.ExtendCopySavedFieldsToDataPool()
		}()
	}
  l.done <- true
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
