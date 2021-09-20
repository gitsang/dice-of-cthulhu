package mooncake

import (
	"cthulhu/pkg/common"
	"math/rand"
	"strings"
	"time"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

func init() {
	SixSideDiceMapInit()
}

var SixSideDiceMap map[int]string

func SixSideDiceMapInit() {
	SixSideDiceMap = make(map[int]string)
	SixSideDiceMap[1] = "⚀"
	SixSideDiceMap[2] = "⚁"
	SixSideDiceMap[3] = "⚂"
	SixSideDiceMap[4] = "⚃"
	SixSideDiceMap[5] = "⚄"
	SixSideDiceMap[6] = "⚅"
	moonCakeDices = NewMoonCakeDices()
}

type MoonCakeDices struct {
}

var moonCakeDices MoonCakeDices

func NewMoonCakeDices() MoonCakeDices {
	return MoonCakeDices{}
}

func (ds *MoonCakeDices) Gamble() (diceStr string, result string) {
	counts := make(map[int]int)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 6; i++ {
		d := rand.Intn(5) + 1
		counts[d]++
		diceStr += SixSideDiceMap[d]
	}

	var hall bool
	for _, cnt := range counts {
		if cnt != 1 {
			hall = false
			break
		}
		hall = true
	}
	if hall {
		result = " 对堂 "
		return
	}

	for k, cnt := range counts {
		if cnt == 1 {
			if k == 4 {
				result += " 一秀 "
			}
		} else if cnt == 2 {
			if k == 4 {
				result += " 二举 "
			}
		} else if cnt == 3 {
			if k == 4 {
				result += " 三红 "
			}
		} else if cnt == 4 {
			if k == 4 {
				result += " 状元 "
			} else {
				result += " 四进 "
			}
		} else if cnt == 5 {
			result += " 状元 "
		}
	}

	log.Debug("gamble result", zap.Reflect("ds", ds))
	return
}

func MoonCakeGamblingInfo(cmd string) *common.Message {
	title := "MoonCake Gambling"
	text := ""

	cmd = strings.TrimPrefix(cmd, ".mooncake")
	cmd = strings.TrimSpace(cmd)
	if cmd == "init" {
		moonCakeDices = NewMoonCakeDices()
		text = "init success"
	}

	return &common.Message{
		MsgType: common.MdType,
		Markdown: common.Markdown{
			Title: title,
			Text:  text,
		},
	}
}

func MoonCakeGambling(usr string) *common.Message {
	title := "MoonCake Gambling"
	text := "# " + usr + " "

	diceStr, result := moonCakeDices.Gamble()
	text += diceStr + "\n" + result

	log.Info("mooncake gambling", zap.String("text", text))
	return &common.Message{
		MsgType: common.MdType,
		Markdown: common.Markdown{
			Title: title,
			Text:  text,
		},
	}
}
