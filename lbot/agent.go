package lbot

import (
  "net/rpc"
  "log"
  "net/http"
  "net"
)
type Agent struct {
	lbot *Lbot
}

func (a *Agent) Listen() error {
  // listen for struct methods
	// arith := new(Arith)
	// rpc.Register(arith)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go http.Serve(l, nil)
  return nil
}
