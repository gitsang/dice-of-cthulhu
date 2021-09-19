package dice

import (
	"cthulhu/pkg/common"
	"math/rand"
	"strconv"
)

func Dice(d int) *common.Message {
	ans := rand.Intn(d)
	content := "d" + strconv.Itoa(d) + " = " + strconv.Itoa(ans)

	return &common.Message{
		MsgType: common.TextType,
		Text: common.Text{
			Content: content,
		},
	}
}
