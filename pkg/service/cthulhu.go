package service

import (
	"bytes"
	"cthulhu/pkg/common"
	"cthulhu/pkg/dice"
	"cthulhu/pkg/mooncake"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

func Help() common.Message {
	title := "Help"
	text := "# Command\n" +
		"---\n" +
		"- `.help`: 查看帮助\n" +
		"- `.d+num`: 骰num面骰子\n" +
		"- `.m`: mooncake gambling\n" +
		"- `.mooncake [option]`: mooncake gambling command\n"

	return common.Message{
		MsgType: common.MdType,
		Markdown: common.Markdown{
			Title: title,
			Text:  text,
		},
	}
}

func parseMessage(reqMsg common.Message) common.Message {
	cmd := reqMsg.Text.Content
	cmd = strings.TrimSpace(cmd)
	log.Info("parse command", zap.String("cmd", cmd))

	if strings.HasPrefix(cmd, ".help") {
		return Help()
	}

	if strings.HasPrefix(cmd, ".d") {
		d, _ := strconv.Atoi(strings.TrimPrefix(cmd, ".d"))
		return dice.Dice(d)
	}

	if strings.HasPrefix(cmd, ".mooncake") {
		return mooncake.MoonCakeGamblingInfo(reqMsg.Text.Content)
	}

	if strings.HasPrefix(cmd, ".m") {
		return mooncake.MoonCakeGambling(reqMsg.SenderNick)
	}

	return common.Message{
		MsgType: common.TextType,
		Text: common.Text{
			Content: "invaild command",
		},
	}
}

func CthulhuHandler(w http.ResponseWriter, r *http.Request) {
	reqMsgJ, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("read body failed", zap.Error(err))
		return
	}

	// parse request message
	var reqMsg common.Message
	err = json.Unmarshal(reqMsgJ, &reqMsg)
	if err != nil {
		log.Error("read body failed", zap.Error(err))
		return
	}
	log.Info("cthulhu recv message", zap.Reflect("reqMsg", reqMsg))
	respMsg := parseMessage(reqMsg)

	// response
	url := reqMsg.SessionWebHook
	respMsgJ, _ := json.Marshal(respMsg)
	result, err := http.Post(url, "application/json", bytes.NewReader(respMsgJ))
	if err != nil {
		log.Error("response failed", zap.Error(err))
		return
	}

	// result
	resJ, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Error("read body failed", zap.Error(err))
		return
	}
	log.Info("response success", zap.ByteString("respJ", respMsgJ), zap.ByteString("result", resJ))
}