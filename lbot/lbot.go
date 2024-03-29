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
	"github.com/pkg/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lbot struct {
	Config         *config.Config
	ctx            context.Context
	mutext         sync.Mutex
	workers        map[string]*worker.Worker
	done           chan bool
	runningAgents  uint64 // todo: remove from here
	changed        chan uint64

  // todo: to move to abstraction
	internalClient *database.MongoClient
}

func NewLbot(ctx context.Context, cfg *config.Config) (*Lbot, error) {
	client, err := database.NewInternalMongoClient(cfg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("Connecting to database failed: %w", err)
	}

	return &Lbot{
		ctx:            ctx,
		Config:         cfg,
		runningAgents:  1,
		changed:        make(chan uint64),
		workers:        map[string]*worker.Worker{},
		internalClient: client,
	}, nil
}

func (l *Lbot) Client() (err error) {
	if l.internalClient == nil {
		return l.internalClient.Disconnect()
	}
	return nil
}

func (l *Lbot) Close() (err error) {
	if l.internalClient == nil {
		return l.internalClient.Disconnect()
	}
	return nil
}

func (l *Lbot) Run() (err error) {
	for _, job := range l.Config.Jobs {
		err = l.internalClient.RunJob(*job)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Lbot) SetCommandState(command *database.Command, state database.CommandState) error {
	command.State = state.String()
	return l.internalClient.SaveCommand(command)
}

func (l *Lbot) SetWorkloadState(workload *database.Workload, state database.WorkloadState) error {
	workload.State = state.String()
	return l.internalClient.SaveWorkload(workload)
}

func (l *Lbot) SetConfig(config *config.Config) {
	l.Config = config
}

func (l *Lbot) StartWorkload(workload *database.Workload) {
	if workload == nil {
		return
	}

	if _, ok := l.workers[workload.Id.String()]; ok {
		log.Println("workload ", workload.Id.String(), " is running")
		return
	}

	l.done = make(chan bool)
	// todo: ping db, before workers init
	// init datapools
	dataPools := make(map[string]schema.DataPool)
	for _, sh := range l.Config.Schemas {
		dataPools[sh.Name] = schema.NewDataPool(sh)
	}

	job := workload.Data
	// // todo: in a parallel depending on type
	func() {
		dataPool := dataPools[job.Schema]

		worker, error := worker.NewWorker(l.ctx, l.Config, &job, dataPool, l.runningAgents)
		if error != nil {
			panic("Worker initialization error")
		}
		fmt.Printf("init worker with job %s\n", job.Name)

		l.mutext.Lock()
		err := l.SetWorkloadState(workload, database.WorkloadStateRunning)
		if err != nil {
			log.Println("error found setting workload done", err)
			return
		}
		l.workers[workload.Id.String()] = worker
		l.mutext.Unlock()
		// todo: fix here, no schema data pool will be nill

		// update: workload state

		defer worker.Close()
		worker.InitMetrics()
		// workaround
		worker.Work(l.changed)
		// worker.Summary()
		worker.ExtendCopySavedFieldsToDataPool()

		l.mutext.Lock()
		err = l.SetWorkloadState(workload, database.WorkloadStateDone)
		if err != nil {
			log.Println("error found setting workload done", err)
		}
		delete(l.workers, workload.Id.String())
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

func (l *Lbot) InitAgent(id primitive.ObjectID, name string) error {
	ct, err := l.internalClient.ClusterTime()
	if err != nil {
		return errors.Wrap(err, "get cluster time")
	}

	agentStatus := database.AgentStatus{
		Id:        id,
		Name:      name,
		CreatedAt: *ct,
	}

	return l.internalClient.AddAgentStatus(agentStatus)
}

func (l *Lbot) AgentHeartBeat(id primitive.ObjectID, name string) error {
	agentStatus := database.AgentStatus{
		Id:   id,
		Name: name,
	}

	return l.internalClient.SaveAgentStatus(agentStatus)
}

func (l *Lbot) HandleWorkload() {
	log.Println("Fetching new workloads")
	// todo: change to commands
	workload, err := l.internalClient.GetNewWorkloads()
	if err != nil {
		return
	}
	log.Println("Fetched new workload with: ", workload.Id.String())

	switch workload.State {
	case database.WorkloadStateCreated.String():

		err := l.SetWorkloadState(workload, database.WorkloadStateToRun)
		if err != nil {
			log.Println("Fetched command with: ", err)
			// if not saved, propably other agent taked
			return
		}

		go l.StartWorkload(workload)
	case database.WorkloadStateToDelete.String():

		// remove awaiting batches
		// change
		// send stop commands
		// todo: add arg workload, if nil stop all
		go l.Cancel()
	}
}

func (l *Lbot) HandleCommand() {
	// todo: change to generic abstraction
	log.Println("Fetching not finished commands")
	// todo: change to commands
	command, err := l.internalClient.GetNextUnFinishedCommand()
	if err != nil {
		return
	}
	log.Println("Fetched command with: ", command.Id.String())

	switch command.State {
	case database.CommandStateCreated.String():
		log.Println("Set command: ", command.Id.String(), " - running")
		if err := l.SetCommandState(command, database.CommandStateRunning); err != nil {
			return
		}
		workloads, err := l.GenerateWorkload(command)
		if err != nil {
			return
		}
		log.Println("Generated new ", len(workloads), " workloads for command: ", command.Id.String())
		if err = l.SaveWorkloads(workloads); err != nil {
			return
		}
		log.Println("Workloads saved successfully")

	case database.CommandStateRunning.String():

		finished, err := l.AreWorkloadsFinished(command)
		if err != nil {
			return
		}

		if finished {
			log.Println("Set command: ", command.Id.String(), " - done")
			if err := l.SetCommandState(command, database.CommandStateDone); err != nil {
				return
			}
		}
		// check if everything is done and I can set command done
	}
}

func (l *Lbot) IsMasterAgent(agentId primitive.ObjectID) (bool, error) {
	return l.internalClient.IsMasterAgent(agentId)
}

// depricated, to be removed
func (l *Lbot) UpdateRunningAgents() error {
	runningAgents, err := l.internalClient.GetAgentWithHeartbeatWithin()
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

func (l *Lbot) GetNextUnFinishedCommand() (*database.Command, error) {
	return l.internalClient.GetNextUnFinishedCommand()
}

func (l *Lbot) AreWorkloadsFinished(command *database.Command) (bool, error) {
	workloads, err := l.internalClient.GetCommandWorkloads(command)
	if err != nil {
		return false, err
	}

	finished := lo.Filter(workloads, func(w *database.Workload, index int) bool {
		return w.State == database.WorkloadStateDone.String() || w.State == database.WorkloadStateError.String()
	})

	return len(workloads) == len(finished) && len(workloads) != 0, nil
}

func (l *Lbot) GenerateWorkload(command *database.Command) ([]*database.Workload, error) {
	ct, err := l.internalClient.ClusterTime()
	if err != nil {
		return nil, errors.Wrap(err, "get cluster time")
	}

	// todo: listener on changed agents, set workloads to error and add new to retry or ??

	// simple approach - each agent gets one workload command
	// command is updated when new running agent occurs

	// when new agent is added
	// master stop running workloads gracefully if we need
	workloads := make([]*database.Workload, 0)

	log.Println("Generating workloads for agents ", l.runningAgents)
	for i := 0; i < int(l.runningAgents); i++ {
		workload := database.Workload{
			Id:        primitive.NewObjectID(),
			CommandId: command.Id,
			Data:      command.Data,
			State:     database.WorkloadStateCreated.String(),
			Version:   primitive.NewObjectID(),
			CreatedAt: *ct,
		}
		workloads = append(workloads, &workload)

	}
	log.Println("Generated workloads: ", len(workloads), workloads)

	return workloads, nil
}

func (l *Lbot) SaveWorkloads(ws []*database.Workload) error {
	var worklods []interface{}
	for _, w := range ws {
		worklods = append(worklods, w)
	}

	return l.internalClient.AddWorkloads(worklods)
}
