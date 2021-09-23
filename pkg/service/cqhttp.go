package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
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

	var msg CQMessage
	err = json.Unmarshal(body, &msg)
	if err != nil {
		log.Error("unmarshal failed", zap.ByteString("body", body), zap.Error(err))
		return
	}

	if strings.Contains(msg.RawMessage, ".d6") {
		res := rand.Intn(6) + 1
		url := fmt.Sprintf("http://127.0.0.1:5700/send_group_msg?group_id=%d&message=%s%d", msg.GroupId, msg.Sender.Nickname, res)
		http.Get(url)
	}

	log.Info("CQ message", zap.Reflect("msg", msg))
}
