package lbot

import (
	"context"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"

	"github.com/kuzxnia/loadbot/lbot/config"
	log "github.com/sirupsen/logrus"
)

type Agent struct {
	ctx    context.Context
	log    *log.Entry
	config *config.Config
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
	rpc.Register(NewStartProcess(a.ctx, a.lbot))
	rpc.Register(NewWatchingProcess(a.ctx, a.lbot))
	rpc.Register(NewStoppingProcess(a.ctx, a.lbot))

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
func (a *Agent) ApplyConfig(configFilePath string) error {
	// todo:
	// check if operation is running
	// lock ?
	// if some operation is running {
	//   return errors.New("")
	// }

	request, err := ParseConfigFile(configFilePath)
	if err != nil {
		return err
	}
	a.lbot.SetConfig(request)

	return nil
}
