package service

import (
	"bytes"
	"cthulhu/pkg/config"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

const (
	gturl   = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	sendurl = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
)

type Text struct {
	Content string `json:"content,omitempty"`
}

type TextCard struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
	BtnText     string `json:"btntext,omitempty"`
}

type SendMsgReq struct {
	ToUser   string    `json:"touser,omitempty"`
	ToParty  string    `json:"toparty,omitempty"`
	ToTag    string    `json:"totag,omitempty"`
	MsgType  string    `json:"msgtype,omitempty"`
	AgentId  int       `json:"agentid,omitempty"`
	Text     *Text     `json:"text,omitempty"`
	TextCard *TextCard `json:"textcard,omitempty"`
	Safe     int       `json:"safe,omitempty"`
}

type GetTokenResp struct {
	Errcode     int    `json:"errcode"`
	AccessToken string `json:"access_token"`
}

func GetToken() (string, error) {
	url := fmt.Sprintf(gturl, config.Cfg.Acl.CorpId, config.Cfg.Acl.CorpSecret)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data GetTokenResp
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	if data.Errcode != 0 {
		return "", errors.New("get token response error " + strconv.Itoa(data.Errcode))
	}

	return data.AccessToken, nil
}

func SendMsg(msg string) error {
	token, err := GetToken()
	if err != nil {
		log.Error("get token failed", zap.Error(err))
		return err
	}

	url := fmt.Sprintf(sendurl, token)
	reqJ, _ := json.Marshal(SendMsgReq{
		ToUser:  "@all",
		MsgType: "text",
		AgentId: config.Cfg.Acl.AgentId,
		Text: &Text{
			Content: msg,
		},
	})

	resp, err := http.Post(url, "application/json", bytes.NewReader(reqJ))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respJ, err := ioutil.ReadAll(resp.Body)

	log.Info("send message success", zap.ByteString("resp", respJ), zap.String("msg", msg))
	return nil
}
