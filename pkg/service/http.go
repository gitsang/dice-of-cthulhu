package service

import (
	"cthulhu/pkg/config"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// redirect
	agent := r.Header.Get("User-Agent")
	if strings.ContainsAny(agent, "CQHttp") {
		CQHandler(w, r)
		return
	}

	// default log
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("can not read body", zap.Error(err))
		return
	}
	log.Debug("receive message",
		zap.String("uri", r.URL.RequestURI()),
		zap.Reflect("header", r.Header),
		zap.ByteString("body", body),
	)
}

func StartHttpServer() {
	http.HandleFunc("/cthulhu", CthulhuHandler)
	http.HandleFunc("/", DefaultHandler)

	listen := strings.Join([]string{config.Cfg.Http.Host, strconv.Itoa(config.Cfg.Http.Port)}, ":")
	log.Info("start httpserver", zap.String("listen", listen))
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		panic(err)
	}
}
