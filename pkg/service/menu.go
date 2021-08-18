package service

import (
	"bytes"
	"cthulhu/pkg/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

const (
	cmurl = "https://qyapi.weixin.qq.com/cgi-bin/menu/create?access_token=%s&agentid=%d"
)

type Button []struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	Url       string `json:"url"`
	SubButton Button `json:"sub_button"`
}

type ButtonReq struct {
	Button Button `json:"button"`
}

func CreateMenu() error {
	token, err := GetToken()
	if err != nil {
		log.Error("get token failed", zap.Error(err))
		return err
	}

	url := fmt.Sprintf(cmurl, token, config.Cfg.Acl.AgentId)

	buttonJ, err := json.Marshal(ButtonReq{
		Button: Button{
			{
				Name: "Goddess Of Dice",
				SubButton: Button{
					{
						Type: "click",
						Name: "d100",
						Key:  ".d100",
					},
					{
						Type: "click",
						Name: "d20",
						Key:  ".d20",
					},
					{
						Type: "click",
						Name: "d10",
						Key:  ".d10",
					},
					{
						Type: "click",
						Name: "d6",
						Key:  ".d6",
					},
				},
			},
		},
	})

	resp, err := http.Post(url, "application/json", bytes.NewReader(buttonJ))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respJ, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Info("create button menu", zap.ByteString("resp", respJ))

	return nil
}
