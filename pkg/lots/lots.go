package lots

import (
	"cthulhu/pkg/common"
	"cthulhu/pkg/config"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
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

	log.Info("lots init success",
		zap.String("lots_path", path),
		zap.Reflect("lotsDatas", lotsDatas),
		zap.Reflect("userDatas", userDatas),
	)
	return nil
}

func GenLots(usr string) *common.Message {
	var (
		err      error
		text     string
		idx      int
		lotsData Pkg
	)

	// init data
	if lotsDatas == nil || userDatas == nil {
		err = Init(config.Cfg.Service.Lots.DataPath)
		if err != nil {
			log.Error("init lots data failed", zap.Error(err))
			return nil
		}
	}
	lotsData = lotsDatas[config.Cfg.Service.Lots.PkgUse]

	// init user data if not exist
	if _, ok := userDatas[usr]; !ok {
		userDatas[usr] = NewUserData()
	}
	usrData := userDatas[usr]

	// already gen
	day := time.Now().Day()
	_, ok := usrData.Lots[day]
	if ok {
		idx = usrData.Lots[day]

		lot := lotsData.Lots[strconv.Itoa(idx)]

		text = "_" + lotsData.Name + "_\n\n" +
			"_今日已签到_\n\n" +
			"---\n\n" +
			"# " + lot.Name + "\n\n"
		contents := strings.Split(lot.Content, " ")
		for _, c := range contents {
			text += "" + c + "\n\n"
		}
	} else {
		rand.Seed(time.Now().UnixNano())
		idx = rand.Intn(len(lotsData.Lots)-1) + 1

		lot := lotsData.Lots[strconv.Itoa(idx)]

		text = "_" + lotsData.Name + "_\n\n" +
			"---\n\n" +
			"# " + lot.Name + "\n\n"
		contents := strings.Split(lot.Content, " ")
		for _, c := range contents {
			text += "" + c + "\n\n"
		}
		usrData.Lots[day] = idx
	}

	log.Info("gen lots",
		zap.String("usr", usr),
		zap.Int("day", day),
		zap.Int("idx", idx),
	)
	return &common.Message{
		MsgType: common.MdType,
		Markdown: common.Markdown{
			Title: "每日一签",
			Text:  text,
		},
	}
}
