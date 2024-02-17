package lbot

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
	"github.com/kuzxnia/loadbot/lbot/proto"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Agent struct {
	ctx        context.Context
	lbot       *Lbot
	grpcServer *grpc.Server
}

func NewAgent(ctx context.Context) *Agent {
	lbot := NewLbot(ctx)

	grpcServer := grpc.NewServer()
	// register commands
	proto.RegisterStartProcessServer(grpcServer, NewStartProcess(ctx, lbot))
	proto.RegisterStopProcessServer(grpcServer, NewStoppingProcess(ctx, lbot))
	proto.RegisterSetConfigProcessServer(grpcServer, NewSetConfigProcess(ctx, lbot))
	proto.RegisterWatchProcessServer(grpcServer, NewWatchingProcess(ctx, lbot))
	proto.RegisterProgressProcessServer(grpcServer, NewProgressProcess(ctx, lbot))
	reflection.Register(grpcServer)

	return &Agent{
		ctx:        ctx,
		lbot:       lbot,
		grpcServer: grpcServer,
	}
}

func (a *Agent) Listen() error {
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(
		stopSignal, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM,
	)
	address := "0.0.0.0:" + a.lbot.config.Agent.Port

	defer func() {
		log.Info("Stopped lbot-agent started on " + address)
		a.grpcServer.GracefulStop()
	}()

	go func() {
		log.Info("Started lbot-agent on " + address)
		tcpListener, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatal("listen error:", err)
			panic(err)
		}
		if err := a.grpcServer.Serve(tcpListener); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	if a.lbot.config.Agent.MetricsExportPort != "" {
		http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
			metrics.WritePrometheus(w, true)
		})
		go func() {
			log.Infof("Started metrics exporter on :%s/metrics", a.lbot.config.Agent.MetricsExportPort)
			http.ListenAndServe(":"+a.lbot.config.Agent.MetricsExportPort, nil)
		}()
	} else if a.lbot.config.Agent.MetricsExportUrl != "" {
		log.Info("Started exporting metrics :", a.lbot.config.Agent.MetricsExportPort)

		metricsLabels := lo.If(
			a.lbot.config.Agent.Name != "",
			fmt.Sprintf(`instance="%s"`, a.lbot.config.Agent.Name),
		).Else("")

		metrics.InitPush(
			a.lbot.config.Agent.MetricsExportUrl,
			10*time.Second, // todo: add interval param
			metricsLabels,
			true,
		)
	}

	<-stopSignal
	fmt.Println("\nReceived stop signal. Exiting.")

	// is this needed?
	_, cancel := context.WithCancel(a.ctx)
	cancel()

	return nil
}

func (a *Agent) ApplyConfig(request *ConfigRequest) error {
	cfg := NewConfig(request)
	a.lbot.SetConfig(cfg)
	return nil
}
