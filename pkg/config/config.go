package config

import (
	log "github.com/gitsang/golog"
	"github.com/jinzhu/configor"
	"go.uber.org/zap"
)

var Version string

const (
	ENVPrefix = "CTHULHU"
)

var Cfg = struct {
	Http struct {
		Host string `default:"0.0.0.0"`
		Port int    `default:"8893"`
	}
	Acl struct {
		CorpId     string
		CorpSecret string
		AgentId    int
		Token      string
		AesKey     string
	}
}{}

func Load(path string) {
	err := configor.New(&configor.Config{
		ENVPrefix: ENVPrefix,
	}).Load(&Cfg, path)
	if err != nil {
		log.Error("load config failed", zap.Error(err))
	}
	log.Info("loag config succes", zap.Reflect("config", Cfg))
}
