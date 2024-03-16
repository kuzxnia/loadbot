package cli

import (
	"context"
	"fmt"

	"github.com/kuzxnia/loadbot/lbot"
	"github.com/kuzxnia/loadbot/lbot/resourcemanager"
	log "github.com/sirupsen/logrus"
)

func InstallResources(request *resourcemanager.InstallRequest) (err error) {
	fmt.Println("🚀 Instalation started")

	orchiestrator := lbot.NewOrchiestrator(context.TODO())

	_, err = orchiestrator.Install(context.TODO(), request)
	if err != nil {
		log.Fatal("arith error:", err)
		return
	}

	fmt.Println("✅ Installation finished sucessfully")

	return nil
}

func UpgradeResources(request *resourcemanager.UpgradeRequest) (err error) {
	fmt.Println("🚀 Instalation started")

	orchiestrator := lbot.NewOrchiestrator(context.TODO())

	_, err = orchiestrator.Upgrade(context.TODO(), request)
	if err != nil {
		log.Fatal("arith error:", err)
		return
	}

	fmt.Println("✅ Installation finished sucessfully")

	return nil
}

func UnInstallResources(request *resourcemanager.UnInstallRequest) (err error) {
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

func ListResources(request *resourcemanager.ListRequest) (err error) {
	fmt.Println("🚀 UnInstalation started")

	orchiestrator := lbot.NewOrchiestrator(context.TODO())

	_, err = orchiestrator.List(context.TODO(), request)
	if err != nil {
		log.Fatal("arith error:", err)
		return
	}

	fmt.Println("✅ UnInstallation finished sucessfully")

	return nil
}
