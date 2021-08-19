package service

import (
	"bytes"
	"cthulhu/pkg/config"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

type AtUser struct {
	DingTalkId string `json:"dingTalkId"`
}

type Text struct {
	Content string `json:"content,omitempty"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Message struct {
	// msg
	MsgId    string   `json:"msgId,omitempty"`
	MsgType  string   `json:"msgtype,omitempty"`
	CreateAt int64    `json:"createAt,omitempty"`
	Text     Text     `json:"text,omitempty"`
	Markdown Markdown `json:"markdown,omitempty"`

	// sender
	SenderId      string   `json:"senderId,omitempty"`
	SenderNick    string   `json:"senderNick,omitempty"`
	SenderCorpId  string   `json:"senderCorpId,omitempty"`
	SenderStaffId string   `json:"senderStaffId,omitempty"`
	AtUsers       []AtUser `json:"atUsers,omitempty"`
	IsAdmin       bool     `json:"isAdmin,omitempty"`
	IsInAtList    bool     `json:"isInAtList,omitempty"`

	// conversation
	ConversationId    string `json:"conversationId,omitempty"`
	ConversationTitle string `json:"conversationTitle,omitempty"`
	ConversationType  string `json:"conversationType,omitempty"`

	// webhook
	SessionWebHook            string `json:"sessionWebHook,omitempty"`
	SessionWebHookExpiredTime int64  `json:"sessionWebHookExpiredTime,omitempty"`

	// bot
	ChatbotUserId string `json:"chatbotUserId,omitempty"`
}

var (
	token  string
	expire int64
)

const (
	TextType = "text"
	MdType   = "markdown"
)

func parseCmd(cmd string) Message {
	if strings.HasPrefix(cmd, ".d") {
		d, _ := strconv.Atoi(strings.TrimPrefix(cmd, ".d"))
		ans := rand.Intn(d)
		content := "d" + strconv.Itoa(d) + " = " + strconv.Itoa(ans)

		return Message{
			MsgType: TextType,
			Text: Text{
				Content: content,
			},
		}
	} else if strings.HasPrefix(cmd, ".help") {
		title := "Help"
		text := "# Command\n" +
			"---\n" +
			"- `.help`: 查看帮助\n" +
			"- `.d+num`: 骰num面骰子\n"

		return Message{
			MsgType: MdType,
			Markdown: Markdown{
				Title: title,
				Text:  text,
			},
		}
	}

	return Message{
		MsgType: TextType,
		Text: Text{
			Content: "invaild command",
		},
	}
}

func CthulhuHandler(w http.ResponseWriter, r *http.Request) {
	raw := r.URL.RawQuery
	bodyB, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("read body failed", zap.Error(err))
		return
	}

	// parse request
	var reqMsg Message
	err = json.Unmarshal(bodyB, &reqMsg)
	if err != nil {
		log.Error("read body failed", zap.Error(err))
		return
	}

	log.Info("cthulhu recv message", zap.String("query", raw), zap.Reflect("reqMsg", reqMsg))

	// parse cmd
	var respMsg Message
	cmd := strings.TrimSpace(reqMsg.Text.Content)
	if strings.HasPrefix(cmd, ".") {
		respMsg = parseCmd(cmd)
	}
	respMsgJ, _ := json.Marshal(respMsg)

	// response
	url := reqMsg.SessionWebHook
	result, err := http.Post(url, "application/json", bytes.NewReader(respMsgJ))
	if err != nil {
		log.Error("response failed", zap.Error(err))
		return
	}

	resultB, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Error("read body failed", zap.Error(err))
		return
	}

	log.Info("response success", zap.ByteString("respJ", respMsgJ), zap.ByteString("result", resultB))
}

func StartHttpServer() {
	http.HandleFunc("/cthulhu", CthulhuHandler)

	listen := strings.Join([]string{config.Cfg.Http.Host, strconv.Itoa(config.Cfg.Http.Port)}, ":")
	log.Info("start httpserver", zap.String("listen", listen))
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		panic(err)
	}
}
