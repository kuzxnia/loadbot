package lbot

import (
	"context"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"

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

func (a *Agent) Listen() error {
	// register driver commands
	// rpc.Register(NewStartProcess(a.ctx, a.lbot))
	rpc.Register(NewWatchingProcess(a.ctx, a.lbot))
	rpc.Register(NewStoppingProcess(a.ctx, a.lbot))
	rpc.Register(NewSetConfigProcess(a.ctx, a.lbot))

	rpc.HandleHTTP()
	agentHost := "0.0.0.0:1234"
	l, err := net.Listen("tcp", agentHost)
	if err != nil {
		a.log.Fatal("listen error:", err)
	}
	a.log.Info("Started lbot-agent on " + agentHost)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go http.Serve(l, nil)

	<-stop
	_, cancel := context.WithCancel(a.ctx)
	cancel()

	a.log.Info("Shuted down lbot-agent")

	return nil
}

func (a *Agent) ListenGRPC() error {
	agentHost := "0.0.0.0:1235"
	l, err := net.Listen("tcp", agentHost)
	if err != nil {
		a.log.Fatal("listen error:", err)
	}

	grpcServer := grpc.NewServer()
	// register commands
	proto.RegisterStartProcessServer(grpcServer, NewStartProcess(a.ctx, a.lbot))

	reflection.Register(grpcServer)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	a.log.Info("Started lbot-agent on " + agentHost)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	<-stop
	_, cancel := context.WithCancel(a.ctx)
	cancel()

	a.log.Info("Shuted down lbot-agent")

	return nil
}

// runned when initializing agent, and after reconfig
func (a *Agent) ApplyConfig(request *ConfigRequest) error {
	// todo:
	// check if operation is running
	// lock ?
	// if some operation is running {
	//   return errors.New("")
	// }

	cfg := NewConfig(request)
	a.lbot.SetConfig(cfg)

	return nil
}
