package common

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
