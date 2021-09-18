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
	"time"

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

func Help() Message {
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

func Dice(d int) Message {
	ans := rand.Intn(d)
	content := "d" + strconv.Itoa(d) + " = " + strconv.Itoa(ans)

	return Message{
		MsgType: TextType,
		Text: Text{
			Content: content,
		},
	}
}

var SixSideDiceMap map[int]string

func init() {
	// init six side dice
	SixSideDiceMap = make(map[int]string)
	SixSideDiceMap[1] = "⚀"
	SixSideDiceMap[2] = "⚁"
	SixSideDiceMap[3] = "⚂"
	SixSideDiceMap[4] = "⚃"
	SixSideDiceMap[5] = "⚄"
	SixSideDiceMap[6] = "⚅"
}

type MoonCakeDices struct {
	Count map[int]int
}

func NewMoonCakeDices() MoonCakeDices {
	return MoonCakeDices{
		Count: make(map[int]int),
	}
}

func (ds MoonCakeDices) Gamble() (diceStr string, result string) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 6; i++ {
		d := rand.Intn(5) + 1
		ds.Count[d]++
		diceStr += SixSideDiceMap[d]
	}

	switch ds.Count[4] {
	case 1:
		result = "yixiu"
	case 2:
		result = "erju"
	case 3:
		result = "sanhong"
	}

	for k, v := range ds.Count {
		if v == 4 {
			if k == 4 {
				result += "zhuangyuan"
			} else {
				result += "sijin"
			}
		}
		if v == 5 {
			result = "zhuangyuan"
		}
	}

	return
}

func MoonCakeGambling(usr string) Message {
	title := "MoonCake Gambling"
	text := "# " + usr + " "

	diceStr, result := NewMoonCakeDices().Gamble()
	text += diceStr + "\n" + result

	log.Info("mooncake gambling", zap.String("text", text))
	return Message{
		MsgType: MdType,
		Markdown: Markdown{
			Title: title,
			Text:  text,
		},
	}
}

func parseMessage(reqMsg Message) Message {
	cmd := reqMsg.Text.Content
	cmd = strings.TrimSpace(cmd)
	log.Info("parse command", zap.String("cmd", cmd))

	if strings.HasPrefix(cmd, ".help") {
		return Help()
	}

	if strings.HasPrefix(cmd, ".d") {
		d, _ := strconv.Atoi(strings.TrimPrefix(cmd, ".d"))
		return Dice(d)
	}

	if strings.HasPrefix(cmd, ".mooncake") {
		return MoonCakeGambling(reqMsg.SenderNick)
	}

	if strings.HasPrefix(cmd, ".m") {
		return MoonCakeGambling(reqMsg.SenderNick)
	}

	return Message{
		MsgType: TextType,
		Text: Text{
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
	var reqMsg Message
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

func StartHttpServer() {
	http.HandleFunc("/cthulhu", CthulhuHandler)

	listen := strings.Join([]string{config.Cfg.Http.Host, strconv.Itoa(config.Cfg.Http.Port)}, ":")
	log.Info("start httpserver", zap.String("listen", listen))
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		panic(err)
	}
}
