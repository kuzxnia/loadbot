package lbot

import (
	"context"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

type Agent struct {
	ctx    context.Context
	log    *log.Entry
	config *Config
	lbot   *Lbot
}

func NewAgent(ctx context.Context, logger *log.Entry) *Agent {
	return &Agent{
		ctx: ctx,
		log: logger,
	}
}

func (a *Agent) Listen() error {
	// register driver commands
	rpc.Register(NewStartProcess(a.ctx))
	rpc.Register(NewWatchingProcess(a.ctx))
	rpc.Register(NewStoppingProcess(a.ctx))

	rpc.HandleHTTP()
	agentHost := "127.0.0.1:1234"
	l, err := net.Listen("tcp", agentHost)
	if err != nil {
		a.log.Fatal("listen error:", err)
	}
	a.log.Info("lbot-agent started on " + agentHost)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go http.Serve(l, nil)

	<-stop

	return nil
}

// runned when initializing agent, and after reconfig
func (a *Agent) ApplyConfig(configFilePath string) {
	// check if operation is running
	// lock ?
}
