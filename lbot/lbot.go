package lbot

import (
	"context"
	"fmt"
	"sync"
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
	mutext        sync.Mutex
	workers       map[string]*worker.Worker
	done          chan bool
	runningAgents uint64 // todo: remove from here
	changed       chan uint64
}

func NewLbot(ctx context.Context) *Lbot {
	return &Lbot{
		ctx:           ctx,
		runningAgents: 1,
		changed:       make(chan uint64),
		workers:       map[string]*worker.Worker{},
	}
}

func (l *Lbot) Run() error {
	client, err := database.NewInternalMongoClient(l.Config.ConnectionString)
	if err != nil {
		return err
	}

	for _, job := range l.Config.Jobs {
		err = client.RunJob(*job)
		if err != nil {
			return err
		}
	}

	return nil
}

// temp
func (l *Lbot) SetCommandRunning(command *database.Command) error {
	client, err := database.NewInternalMongoClient(l.Config.ConnectionString)
	if err != nil {
		return err
	}
	return client.SetCommandRunning(command)
}

func (l *Lbot) SetCommandDone(command *database.Command) error {
	client, err := database.NewInternalMongoClient(l.Config.ConnectionString)
	if err != nil {
		return err
	}
	return client.SetCommandDone(command)
}

func (l *Lbot) StartWorkload(command *database.Command) {
	if command == nil {
		return
	}

	if _, ok := l.workers[command.Id.String()]; ok {
		log.Println("Command ", command.Id.String(), " is running")
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

		l.mutext.Lock()
		err := l.SetCommandRunning(command)
		if err != nil {
			log.Println("error found setting command done", err)
			return
		}
		l.workers[command.Id.String()] = worker
		l.mutext.Unlock()
		// todo: fix here, no schema data pool will be nill

		// update: command state

		defer worker.Close()
		worker.InitMetrics()
		// workaround
		worker.Work(l.changed)
		// worker.Summary()
		worker.ExtendCopySavedFieldsToDataPool()

		l.mutext.Lock()
		err = l.SetCommandDone(command)
		if err != nil {
			log.Println("error found setting command done", err)
		}
		delete(l.workers, command.Id.String())
		l.mutext.Unlock()
	}()
	l.done <- true
}

func (l *Lbot) Cancel() error {
	for _, worker := range l.workers {
		worker.Cancel()
	}
	l.workers = map[string]*worker.Worker{}

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
	log.Println("Fetched new command with type: ", command.Type)

	switch command.Type {
	case database.CommandTypeStartWorkload.String():
		go l.StartWorkload(command)
	case database.CommandTypeStopWorkload.String():
		// remove awaiting batches
		// change
		// send stop commands
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
