package main

import (
	"cthulhu/pkg/config"
	"cthulhu/pkg/service"
	"flag"
	"os"
	"os/signal"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

var confPath = flag.String("c", "../conf/cthulhu.yml", "config file")
var logPath = flag.String("l", "../log/cthulhu.log", "log file")

func Init() {
	flag.Parse()
	err := config.Load(*confPath)
	if *logPath != "" {
		config.Cfg.Log.Path = *logPath
	}
	log.InitLogger(
		log.WithLogFile(config.Cfg.Log.Path),
		log.WithLogLevel(config.Cfg.Log.Level),
		log.WithEncoderType(log.EncoderTypeConsole),
		log.WithDisplayFuncEnable(true),
	)
	if err != nil {
		log.Error("load config failed", zap.Error(err))
	}
	log.Info("loag config succes", zap.Reflect("config", config.Cfg))
}

func main() {
	// init
	Init()

	// start
	go service.StartHttpServer()

	// wait
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	select {
	case sig, _ := <-sigChan:
		log.Warn("capture system signal", zap.Any("sig", sig))
	}
}
