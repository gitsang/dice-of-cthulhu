package service

import (
	"cthulhu/pkg/lots"
	"encoding/json"
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

/*
{
	"anonymous": null,
	"font": 0,
	"group_id": 892187734,
	"message": "hello_world",
	"message_id": 975152535,
	"message_seq": 628,
	"message_type": "group",
	"post_type": "message",
	"raw_message": "hello_world",
	"self_id": 1131986664,
	"sender": {
		"age": 0,
		"area": "",
		"card": "",
		"level": "",
		"nickname": "Sangria",
		"role": "owner",
		"sex": "unknown",
		"title": "",
		"user_id": 1203869957
	},
	"sub_type": "normal",
	"time": 1632379682,
	"user_id": 1203869957
}
*/

type CQSender struct {
	Nickname string `json:"nickname"`
	UserId   int64  `json:"user_id"`
}

type CQMessage struct {
	GroupId    int64    `json:"group_id"`
	RawMessage string   `json:"raw_message"`
	Sender     CQSender `json:"sender"`
}

func CQHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("can not read body", zap.Error(err))
		return
	}

	var reqMsg CQMessage
	err = json.Unmarshal(body, &reqMsg)
	if err != nil {
		log.Error("unmarshal failed", zap.ByteString("body", body), zap.Error(err))
		return
	}

	if reqMsg.GroupId == 0 {
		log.Debug("not group message, ignore", zap.ByteString("body", body))
		return
	}

	var respMsg string
	//var title string
	//var content string
	if strings.Contains(reqMsg.RawMessage, ".d6") {
		res := rand.Intn(6) + 1
		respMsg = reqMsg.Sender.Nickname + "\n" +
			"dice result: " + strconv.Itoa(res)

		//title = "dice"
		//content = "dice result: " + strconv.Itoa(res)
	} else if strings.Contains(reqMsg.RawMessage, ".lots") {
		respMsg = lots.GenLots(reqMsg.Sender.Nickname).Markdown.Text

		//title = "dice"
		//content = lots.GenLots(reqMsg.Sender.Nickname).Markdown.Text
	} else {
		log.Debug("not command, ignore", zap.String("RawMessage", reqMsg.RawMessage))
		return
	}

	//respMsg = fmt.Sprintf(`[CQ:xml,data=<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
	//	<msg serviceID="1">
	//		<item>
	//			<title>%s</title>
	//			<summary>%s</summary>
	//			<summary>%s</summary>
	//		</item>
	//	</msg>
	//]`, title, content, "test msg")

	// response
	param := url.Values{}
	param.Add("group_id", fmt.Sprintf("%d", reqMsg.GroupId))
	param.Add("auto_escape", "false")
	param.Add("message", respMsg)

	u := "http://127.0.0.1:5700/send_group_msg"
	u = strings.Join([]string{u, param.Encode()}, "?")

	http.Get(u)
	log.Info("cq message", zap.Reflect("req", reqMsg), zap.String("resp", u))
}
