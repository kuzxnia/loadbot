package agent

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/kuzxnia/loadbot/lbot"
	"github.com/kuzxnia/loadbot/lbot/config"
	"github.com/kuzxnia/loadbot/lbot/proto"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Agent struct {
	ctx        context.Context
	lbot       *lbot.Lbot
	grpcServer *grpc.Server
}

func NewAgent(ctx context.Context, loadbot *lbot.Lbot) *Agent {
	// loadbot := lbot.NewLbot(ctx)

	grpcServer := grpc.NewServer()
	// register commands
	// worker commands
	proto.RegisterStartProcessServer(grpcServer, lbot.NewStartProcess(ctx, loadbot))
	proto.RegisterStopProcessServer(grpcServer, lbot.NewStoppingProcess(ctx, loadbot))
	proto.RegisterSetConfigProcessServer(grpcServer, lbot.NewSetConfigProcess(ctx, loadbot))
	proto.RegisterWatchProcessServer(grpcServer, lbot.NewWatchingProcess(ctx, loadbot))
	proto.RegisterProgressProcessServer(grpcServer, lbot.NewProgressProcess(ctx, loadbot))

	reflection.Register(grpcServer)

	return &Agent{
		ctx:        ctx,
		lbot:       loadbot,
		grpcServer: grpcServer,
	}
}

func (a *Agent) Start() error {
	go a.Listen()
	go a.Metrics()
	go a.Heartbeat()

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
func (a *Agent) Listen() error {
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

func (a *Agent) Heartbeat() error {
	ticker := time.NewTicker(config.AgentsHeartbeatInterval)
	defer ticker.Stop()

	// case agent restarted during update should start from place where ends:
	// generate version at start and check on start if where are data in agent status collection
	// if so, read log and get status

	// if new, and command is running

	for range ticker.C {
		// check how many agents are reachable
		// check only on oldest node, no need for thread race

		// flush workload progress to command
		// ?? colleciton lbotWorkloadState ? and state inserted per agent per workload

		err := a.lbot.RefreshAgentStatus(a.lbot.Config.Agent.Name)
		if err != nil {
			log.Error("agent status failed", err)
		}

		err = a.lbot.UpdateRunningAgents()
		if err != nil {
			log.Error("agent list failed", err)
		}
	}

	return nil
}

// agent config should be removed from base config and later
// zarzadzanie workloadem przez
