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

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dingtalkoauth2_1_0 "github.com/alibabacloud-go/dingtalk/oauth2_1_0"
	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

type AtUser struct {
	DingTalkId string `json:"dingTalkId"`
}

type Text struct {
	Content string `json:"content,omitempty"`
}

type Message struct {
	// msg
	MsgId    string `json:"msgId,omitempty"`
	MsgType  string `json:"msgtype,omitempty"`
	CreateAt int64  `json:"createAt,omitempty"`
	Text     Text   `json:"text,omitempty"`

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

func parseCmd(cmd string) string {
	if strings.HasPrefix(cmd, ".d") {
		d, _ := strconv.Atoi(strings.TrimPrefix(cmd, ".d"))
		ans := rand.Intn(d)
		result := "d" + strconv.Itoa(d) + " = " + strconv.Itoa(ans)
		return result
	}

	return ""
}

func GetToken() (string, error) {
	if token != "" || time.Now().Unix() < expire {
		return token, nil
	}

	var (
		err      error
		client   *dingtalkoauth2_1_0.Client
		protocol = "https"
		region   = "central"
	)

	client, err = dingtalkoauth2_1_0.NewClient(&openapi.Config{
		Protocol: &protocol,
		RegionId: &region,
	})
	if err != nil {
		return "", err
	}

	getAccessTokenRequest := &dingtalkoauth2_1_0.GetAccessTokenRequest{
		AppKey:    &config.Cfg.Acl.AppKey,
		AppSecret: &config.Cfg.Acl.AppSecret,
	}

	resp, err := client.GetAccessToken(getAccessTokenRequest)
	if err != nil {
		return "", err
	}

	token = *resp.Body.AccessToken
	expire = *resp.Body.ExpireIn
	return token, nil
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
	var content string
	cmd := strings.TrimSpace(reqMsg.Text.Content)
	if strings.HasPrefix(cmd, ".") {
		content = parseCmd(cmd)
	}

	// response
	respMsg := Message{
		MsgType: "text",
		Text: Text{
			Content: content,
		},
	}
	respMsgJ, _ := json.Marshal(respMsg)

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
