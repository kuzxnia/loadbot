package agent

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/fsnotify/fsnotify"
	"github.com/kuzxnia/loadbot/lbot"
	"github.com/kuzxnia/loadbot/lbot/config"
	"github.com/kuzxnia/loadbot/lbot/proto"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//go:generate stringer -type=AgentState -trimprefix=AgentState
type AgentState int

const (
	AgentStateLeader AgentState = iota
	AgentStateFollower
)

type Agent struct {
	id           primitive.ObjectID
	ctx          context.Context
	lbot         *lbot.Lbot
	grpcServer   *grpc.Server
	state        AgentState
	stateChange  *sync.Cond
	configChange *fsnotify.Watcher
}

func NewAgent(ctx context.Context, loadbot *lbot.Lbot) *Agent {
	grpcServer := grpc.NewServer()
	// register commands
	proto.RegisterStartProcessServer(grpcServer, lbot.NewStartProcess(ctx, loadbot))
	proto.RegisterStopProcessServer(grpcServer, lbot.NewStoppingProcess(ctx, loadbot))
	proto.RegisterSetConfigProcessServer(grpcServer, lbot.NewSetConfigProcess(ctx, loadbot))
	proto.RegisterWatchProcessServer(grpcServer, lbot.NewWatchingProcess(ctx, loadbot))
	proto.RegisterProgressProcessServer(grpcServer, lbot.NewProgressProcess(ctx, loadbot))

	reflection.Register(grpcServer)

	return &Agent{
		id:          primitive.NewObjectID(),
		ctx:         ctx,
		lbot:        loadbot,
		grpcServer:  grpcServer,
		state:       AgentStateFollower,
		stateChange: sync.NewCond(&sync.Mutex{}),
	}
}

func (a *Agent) Start() error {
	defer func() {
		if a.configChange != nil {
			a.configChange.Close()
		}
	}()

	go a.ServeGrpc()
	go a.Metrics()
	go a.Heartbeat()
	go a.Listen()

	if err := a.lbot.InitAgent(a.id, a.lbot.Config.Agent.Name); err != nil {
	} else {
		log.Info("lbot-agent initialized successfuly")
	}

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(
		stopSignal, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM,
	)
	<-stopSignal
	fmt.Println("\nReceived stop signal. Exiting.")

	// is this needed?
	_, cancel := context.WithCancel(a.ctx)
	cancel()

	return nil
}

// właściwie to nie ma potrzeby nasłuchiwać na grpc dla każdego followera
func (a *Agent) ServeGrpc() error {
	address := "0.0.0.0:" + a.lbot.Config.Agent.Port

	defer func() {
		log.Info("Stopped lbot-agent started on " + address)
		a.grpcServer.GracefulStop()
	}()

	log.Info("Started lbot-agent on " + address)
	tcpListener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("listen error:", err)
		panic(err)
	}
	if err := a.grpcServer.Serve(tcpListener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	return nil
}

// remove from here
func (a *Agent) Metrics() error {
	if lo.IsNotEmpty(a.lbot.Config.Agent.MetricsExportPort) {
		http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
			metrics.WritePrometheus(w, true)
		})
		log.Infof("Started metrics exporter on :%s/metrics", a.lbot.Config.Agent.MetricsExportPort)
		http.ListenAndServe(":"+a.lbot.Config.Agent.MetricsExportPort, nil)
	} else if lo.IsNotEmpty(a.lbot.Config.Agent.MetricsExportUrl) {
		log.Info("Started exporting metrics to ", a.lbot.Config.Agent.MetricsExportUrl)

		metricsLabels := lo.If(
			a.lbot.Config.Agent.Name != "",
			fmt.Sprintf(`instance="%s"`, a.lbot.Config.Agent.Name),
		).Else("")

		metrics.InitPush(
			a.lbot.Config.Agent.MetricsExportUrl,
			10*time.Second, // todo: add interval param
			metricsLabels,
			true,
		)
	}
	return nil
}

// remove from here
func (a *Agent) ApplyConfig(request *lbot.ConfigRequest) error {
	cfg := lbot.NewConfig(request)
	a.lbot.SetConfig(cfg)
	return nil
}

func (a *Agent) ApplyConfigFromFile(path string) error {
	request, err := lbot.ParseConfigFile(path)
	if err != nil {
		return err
	}
	cfg := lbot.NewConfig(request)
	a.lbot.SetConfig(cfg)
	return nil
}

func (a *Agent) Heartbeat() error {
	ticker := time.NewTicker(config.AgentsHeartbeatInterval)
	defer ticker.Stop()

	// case agent restarted during update should start from place where ends:
	// generate version at start and check on start if where are data in agent status collection
	// if so, read log and get status

	// if new, and command is running

	// flow:
	// agent have two modes:
	// worker or coordinator

	// worker requirements
	// - do the work (process commands)
	// - send heartbeats about his state
	// - check if is master
	// store his internal state somewhere (to know where did he finished when crushed) ? - not for this version

	// coordinator / master: when i can be elected as master? when I'm oldest node in set
	// - publish work batch requests to queue
	// 1. adds command to db "command", 'config' - because config could change later
	// - workload uuid - require to distinguish workload
	// workload state - or is finished

	// default batch size is rps * 10(approx 10s of work), or 10s(timed job), or just generate batch for every agent (job without limits)
	// - get agent states (needed to know to wchich agent we should to create batches)
	// - send heartbeats about his state

	// master ticker should run more often, hint ticker.reset

	for range ticker.C {
		err := a.lbot.AgentHeartBeat(a.id, a.lbot.Config.Agent.Name)
		if err != nil {
			log.Error("agent status failed", err)
		}

	}

	return nil
}

// add command struct

// start command
//
// divide workload into groups?
// when adding new agent, stop worklflow, make snapshot, recalculate jobs, and then start
// workload command is root command, each agent creates own workload version, or is part of this command(inside as list item)
//

func (a *Agent) Listen() error {
	// todo:
	// 1. commands - to handle on master
	// 2. subcommands - for worker
	// 3. agents - handle is master

	// a.lbot.HandleCommands()
	// todo: change to stream - subscribe, only new commands

	// each worker subscribe on subcommands
	ticker := time.NewTicker(config.AgentsHeartbeatInterval)
	defer ticker.Stop()

	for range ticker.C {
		// check for running commands
		a.lbot.HandleWorkload()

		// todo: name should be on agent
		// todo: need to notify listen goroutine about data change
		shouldIBeLeader, err := a.lbot.IsMasterAgent(a.id)
		amILeader := a.state == AgentStateLeader
		if shouldIBeLeader != amILeader {
			a.stateChange.L.Lock()
			a.state = lo.If(shouldIBeLeader == true, AgentStateLeader).Else(AgentStateFollower)
			log.Println("new state: ", a.state)
			a.stateChange.L.Unlock()
			a.stateChange.Signal()
			// todo: log
		}

		err = a.lbot.UpdateRunningAgents()
		if err != nil {
			log.Error("agent list failed", err)
		}

		if a.state == AgentStateLeader {
			a.lbot.HandleCommand()
		}
	}

	return nil
}
