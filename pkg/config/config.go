package config

import (
	"github.com/jinzhu/configor"
)

var Version string

const (
	ENVPrefix = "CTHULHU"
)

var Cfg = struct {
	Log struct {
		Level string `default:"Info"`
		Path  string `default:"../log/cthulhu.log"`
	}
	Http struct {
		Host string `default:"0.0.0.0"`
		Port int    `default:"8893"`
	}
	Acl struct {
		CorpId    string
		Token     string
		AgentId   int
		AppKey    string
		AppSecret string
	}
	Service struct {
		Lots struct {
			DataPath string
		}
	}
}{}

func Load(path string) error {
	return configor.New(&configor.Config{
		ENVPrefix: ENVPrefix,
	}).Load(&Cfg, path)
}
