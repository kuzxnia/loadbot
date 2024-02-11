package lbot

import (
	"context"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/VictoriaMetrics/metrics"
	"github.com/kuzxnia/loadbot/lbot/config"
	"github.com/kuzxnia/loadbot/lbot/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Agent struct {
	ctx    context.Context
	log    *log.Entry
	config *config.Config
	lbot   *Lbot
}

func NewAgent(ctx context.Context, logger *log.Entry) *Agent {
	return &Agent{
		ctx:  ctx,
		log:  logger,
		lbot: NewLbot(ctx),
	}
}

func (a *Agent) Listen(port string) error {
	address := "0.0.0.0:" + port

	l, err := net.Listen("tcp", address)
	if err != nil {
		a.log.Fatal("listen error:", err)
	}

	grpcServer := grpc.NewServer()
	// register commands
	proto.RegisterStartProcessServer(grpcServer, NewStartProcess(a.ctx, a.lbot))
	proto.RegisterStopProcessServer(grpcServer, NewStoppingProcess(a.ctx, a.lbot))
	proto.RegisterSetConfigProcessServer(grpcServer, NewSetConfigProcess(a.ctx, a.lbot))
	proto.RegisterWatchProcessServer(grpcServer, NewWatchingProcess(a.ctx, a.lbot))
	proto.RegisterProgressProcessServer(grpcServer, NewProgressProcess(a.ctx, a.lbot))

	reflection.Register(grpcServer)

	a.log.Info("Started lbot-agent on " + address)
	stop := make(chan os.Signal, 1)
	signal.Notify(
		stop, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM,
	)

	go func() {
		if err := grpcServer.Serve(l); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	// register metrics
	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, true)
	})
	go func() {
		a.log.Info("Started lbot-agent metrics on :8090")
		http.ListenAndServe(":8090", nil)
	}()

	<-stop
	grpcServer.GracefulStop()

	// is this needed?
	_, cancel := context.WithCancel(a.ctx)
	cancel()

	a.log.Info("Shuted down lbot-agent")

	return nil
}

// runned when initializing agent, and after reconfig
func (a *Agent) ApplyConfig(request *ConfigRequest) error {
	// todo:
	// check if operation is running
	// lock ? or apply config and restart
	// if some operation is running {
	//   return errors.New("")
	// }

	cfg := NewConfig(request)
	a.lbot.SetConfig(cfg)

	return nil
}
