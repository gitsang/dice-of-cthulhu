package service

import (
	"cthulhu/pkg/config"
	"cthulhu/pkg/crypt/jsonc"
	"cthulhu/pkg/crypt/xmlc"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

func parseCmd(cmd string) {
	if strings.HasPrefix(cmd, ".d") {
		d, _ := strconv.Atoi(strings.TrimPrefix(cmd, ".d"))
		ans := rand.Intn(d)
		SendMsg("d" + strconv.Itoa(d) + " = " + strconv.Itoa(ans))
	}
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	var (
		u         = r.URL.Query()
		sig       = u.Get("msg_signature")
		ts        = u.Get("timestamp")
		nonce     = u.Get("nonce")
		echo, _   = url.QueryUnescape(u.Get("echostr"))
		respB     []byte
		logFields = []zap.Field{
			zap.String("msg_signature", sig),
			zap.String("timestamp", ts),
			zap.String("nonce", nonce),
			zap.String("echostr", echo),
		}
	)
	log.Info("verify url", logFields...)
	defer log.Info("verify url success", append(logFields, zap.ByteString("resp", respB))...)

	wxcpt := jsonc.NewWXBizMsgCrypt(config.Cfg.Acl.Token, config.Cfg.Acl.AesKey, config.Cfg.Acl.CorpId, jsonc.JsonType)
	respB, cptErr := wxcpt.VerifyURL(sig, ts, nonce, echo)
	if cptErr != nil {
		log.Error("verify url failed", zap.Reflect("error", cptErr))
		return
	}

	_, _ = w.Write(respB)
}

type MsgContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   uint32 `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	Event        string `xml:"Event"`
	EventKey     string `xml:"EventKey"`
	Msgid        uint64 `xml:"MsgId"`
	Agentid      uint32 `xml:"AgentId"`
}

func MsgHandler(w http.ResponseWriter, r *http.Request) {
	var (
		u         = r.URL.Query()
		sig       = u.Get("msg_signature")
		ts        = u.Get("timestamp")
		nonce     = u.Get("nonce")
		respB     []byte
		logFields = []zap.Field{
			zap.String("msg_signature", sig),
			zap.String("timestamp", ts),
			zap.String("nonce", nonce),
		}
	)
	log.Info("get new message", logFields...)
	defer log.Info("message parse end", append(logFields, zap.ByteString("resp", respB))...)

	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("read body failed", zap.Error(err))
		return
	}
	log.Info("read body success", zap.ByteString("body", body))

	// decrypt
	wxcpt := xmlc.NewWXBizMsgCrypt(config.Cfg.Acl.Token, config.Cfg.Acl.AesKey, config.Cfg.Acl.CorpId, xmlc.XmlType)
	msg, cptErr := wxcpt.DecryptMsg(sig, ts, nonce, body)
	if cptErr != nil {
		log.Error("crypt failed", zap.Reflect("error", cptErr))
		return
	}
	log.Info("decrypt mseeage success", zap.ByteString("msg", msg))

	// parse
	var msgContent MsgContent
	err = xml.Unmarshal(msg, &msgContent)
	if nil != err {
		log.Error("xml unmarshal failed", zap.Error(err))
		return
	}
	log.Info("xml unmarshal success", zap.Reflect("msg", msgContent))

	// parse cmd/msg
	if strings.HasPrefix(msgContent.Content, ".") {
		_ = SendMsg(msgContent.FromUsername + "(cmd): " + msgContent.Content)
		parseCmd(msgContent.Content)
	} else if strings.HasPrefix(msgContent.EventKey, ".") {
		_ = SendMsg(msgContent.FromUsername + "(cmd): " + msgContent.EventKey)
		parseCmd(msgContent.EventKey)
	} else {
		_ = SendMsg(msgContent.FromUsername + "(msg): " + msgContent.Content)
	}

	// resp
	ToUsername := msgContent.ToUsername
	msgContent.ToUsername = msgContent.FromUsername
	msgContent.FromUsername = ToUsername
	replayJson, err := json.Marshal(&msgContent)

	encryptMsg, cryptErr := wxcpt.EncryptMsg(string(replayJson), "1409659589", "1409659589")
	if nil != cryptErr {
		fmt.Println("DecryptMsg fail", cryptErr)
	}

	_, _ = w.Write(encryptMsg)
}

func StartHttpServer() {
	http.HandleFunc("/cthulhu", func(w http.ResponseWriter, r *http.Request) {
		raw := r.URL.RawQuery

		log.Info("cthulhu handler start")
		defer log.Info("cthulhu handler end", zap.String("raw", raw))

		if strings.Contains(raw, "echostr") {
			VerifyHandler(w, r)
		} else {
			MsgHandler(w, r)
		}
	})

	listen := strings.Join([]string{config.Cfg.Http.Host, strconv.Itoa(config.Cfg.Http.Port)}, ":")
	log.Info("start httpserver", zap.String("listen", listen))
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		panic(err)
	}
}
