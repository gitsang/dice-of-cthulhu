package main

import (
	"cthulhu/pkg/config"
	"cthulhu/pkg/service"
	"os"
	"os/signal"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

func main() {
	config.Load("../conf/cthulhu.yml")

	go service.StartHttpServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	select {
	case sig, _ := <-sigChan:
		log.Warn("capture system signal", zap.Any("sig", sig))
	}
}
