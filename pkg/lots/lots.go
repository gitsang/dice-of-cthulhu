package lots

import (
	"cthulhu/pkg/common"
	"cthulhu/pkg/config"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

type Lot struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Pkg struct {
	Name string                    `json:"name"`
	Lots map[string] /* idx */ Lot `json:"lots"`
}

type UserData struct {
	Lots map[int] /* day */ int /* lots idx */
}

func NewUserData() UserData {
	return UserData{
		Lots: make(map[int]int),
	}
}

var lotsDatas map[string] /* idx */ Pkg
var userDatas map[string] /*name*/ UserData

func Init(path string) error {
	dataJ, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(dataJ, &lotsDatas)
	if err != nil {
		return err
	}

	userDatas = make(map[string]UserData)

	return nil
}

func GenLots(usr string) *common.Message {
	var (
		err      error
		text     string
		lotsData = lotsDatas["1"]
	)

	// init data
	if lotsDatas == nil || userDatas == nil {
		err = Init(config.Cfg.Service.Lots.DataPath)
		if err != nil {
			log.Error("init lots data failed", zap.Error(err))
			return nil
		}
	}

	// init user data if not exist
	if _, ok := userDatas[usr]; !ok {
		userDatas[usr] = NewUserData()
	}
	usrData := userDatas[usr]

	// already gen
	day := time.Now().Day()
	if idx, ok := usrData.Lots[day]; ok {
		text = "already gen\n" +
			lotsData.Lots[strconv.Itoa(idx)].Content
	} else {
		rand.Seed(time.Now().UnixNano())
		idx = rand.Intn(len(lotsData.Lots))
		text = lotsData.Lots[strconv.Itoa(idx)].Content
		usrData.Lots[day] = idx
	}

	return &common.Message{
		MsgType: common.MdType,
		Markdown: common.Markdown{
			Title: "每日一签",
			Text:  text,
		},
	}
}
