package service

import (
	"cthulhu/pkg/config"
	"net/http"
	"strconv"
	"strings"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

func StartHttpServer() {
	http.HandleFunc("/cthulhu", CthulhuHandler)

	listen := strings.Join([]string{config.Cfg.Http.Host, strconv.Itoa(config.Cfg.Http.Port)}, ":")
	log.Info("start httpserver", zap.String("listen", listen))
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		panic(err)
	}
}
