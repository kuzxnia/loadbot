package cli

import (
	"context"
	"fmt"

	"github.com/kuzxnia/loadbot/lbot"
	"github.com/kuzxnia/loadbot/lbot/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func InstallResources(conn grpc.ClientConnInterface, request *proto.InstallRequest) (err error) {
	fmt.Println("🚀 Instalation started")

	// client := proto.NewOrchistratorServiceClient(conn)
  orchiestrator := lbot.NewOrchiestrator(context.TODO())

	_, err = orchiestrator.Install(context.TODO(), request)
	if err != nil {
		log.Fatal("arith error:", err)
		return
	}

	fmt.Println("✅ Installation finished sucessfully")

	return nil
}

func UnInstallResources(conn grpc.ClientConnInterface, request *proto.UnInstallRequest) (err error) {
	fmt.Println("🚀 UnInstalation started")

  orchiestrator := lbot.NewOrchiestrator(context.TODO())

	_, err = orchiestrator.UnInstall(context.TODO(), request)
	if err != nil {
		log.Fatal("arith error:", err)
		return
	}

	fmt.Println("✅ UnInstallation finished sucessfully")

	return nil
}
