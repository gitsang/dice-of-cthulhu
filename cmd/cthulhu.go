package main

import (
	"cthulhu/pkg/config"
	"cthulhu/pkg/service"
	"os"
	"os/signal"
	"time"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

func main() {
	config.Load("../conf/cthulhu.yml")

	service.SendMsg("start:" + time.Now().String())
	_ = service.CreateMenu()
	go service.StartHttpServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	select {
	case sig, _ := <-sigChan:
		log.Warn("capture system signal", zap.Any("sig", sig))
		service.SendMsg("end:" + time.Now().String())
	}
}
