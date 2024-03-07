package lbot

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/kuzxnia/loadbot/lbot/config"
	"github.com/kuzxnia/loadbot/lbot/database"
	"github.com/kuzxnia/loadbot/lbot/schema"
	"github.com/kuzxnia/loadbot/lbot/worker"
	log "github.com/sirupsen/logrus"
)

type Lbot struct {
	Config        *config.Config
	ctx           context.Context
	workers       []*worker.Worker
	done          chan bool
	runningAgents uint64
	changed       chan uint64
}

func NewLbot(ctx context.Context) *Lbot {
	return &Lbot{
		ctx:           ctx,
		runningAgents: 1,
		changed:       make(chan uint64),
	}
}

func (l *Lbot) Run() error {
	client, err := database.NewInternalMongoClient(l.Config.ConnectionString)
	if err != nil {
		return err
	}

	for _, job := range l.Config.Jobs {
		client.RunJob(*job)
	}

	return nil
}

func (l *Lbot) StartWorkload(command *database.Command) {
	if command == nil {
		return
	}
	l.done = make(chan bool)
	// todo: ping db, before workers init
	// init datapools
	dataPools := make(map[string]schema.DataPool)
	for _, sh := range l.Config.Schemas {
		dataPools[sh.Name] = schema.NewDataPool(sh)
	}

	job := command.Data
	// // todo: in a parallel depending on type
	func() {
		dataPool := dataPools[job.Schema]

		worker, error := worker.NewWorker(l.ctx, l.Config, &job, dataPool, l.runningAgents)
		if error != nil {
			panic("Worker initialization error")
		}
		fmt.Printf("init worker with job %s\n", job.Name)
		l.workers = append(l.workers, worker)
		// todo: fix here, no schema data pool will be nill
		defer worker.Close()
		worker.InitMetrics()
		// workaround
		worker.Work(l.changed)
		// worker.Summary()
		worker.ExtendCopySavedFieldsToDataPool()
	}()
	l.done <- true
}

func (l *Lbot) Cancel() error {
	for _, worker := range l.workers {
		worker.Cancel()
	}
	l.workers = nil

	return nil
}

func (l *Lbot) RefreshAgentStatus(name string) error {
	// todo: change to generic abstraction
	client, err := database.NewInternalMongoClient(l.Config.ConnectionString)
	if err != nil {
		return err
	}

	agentStatus := database.AgentStatus{
		Name: name,
	}

	return client.SetAgentStatus(agentStatus)
}

func (l *Lbot) HandleCommands() {
	// todo: change to generic abstraction
	client, err := database.NewInternalMongoClient(l.Config.ConnectionString)
	if err != nil {
		return
	}

	// todo: change to commands
	command, err := client.GetNewCommand()
	if err != nil {
		return
	}
	switch command.Type {
	case database.CommandTypeStartWorkload.String():
		go l.StartWorkload(command)
	case database.CommandTypeStopWorkload.String():
		go l.Cancel()
	}
}

func (l *Lbot) IsMasterAgent(name string) (bool, error) {
	client, err := database.NewInternalMongoClient(l.Config.ConnectionString)
	if err != nil {
		return false, err
	}
	return client.IsMasterAgent(name)
}

// depricated, to be removed
func (l *Lbot) UpdateRunningAgents() error {
	client, err := database.NewInternalMongoClient(l.Config.ConnectionString)
	if err != nil {
		return err
	}

	runningAgents, err := client.GetAgentWithHeartbeatWithin()
	if err != nil {
		return err
	}
	if runningAgents != l.runningAgents {
		log.Info("New running agents value ", runningAgents)
		atomic.StoreUint64(&l.runningAgents, runningAgents)
		select {
		case l.changed <- runningAgents:
			log.Info("Workers notified, new running agents value ", runningAgents)
		default:
		}
	}

	return nil
}

func (l *Lbot) GetNotFinishedCommands() ([]*database.Command, error) {
	return nil, nil
}

func (l *Lbot) GenerateWorkerSubCommands() ([]*database.SubCommand, error) {
	return nil, nil
}

func (l *Lbot) SetCommandFinished(command *database.Command) error {
	return nil
}

func (l *Lbot) InsertWorkerSubCommands(subCommands []*database.SubCommand) error {
	return nil
}

func (l *Lbot) SetConfig(config *config.Config) {
	l.Config = config
}
