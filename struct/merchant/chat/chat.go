package MerchantChatStruct

type ChatData struct {
	ChatId       int    `json:"chatId"`
	ChatToken    string `json:"chatToken"`
	Open         bool   `json:"open"`
	CreateDate   string `json:"createDate"`
	LastTalkDate string `json:"lastTalkDate"`
}

type ChatMessage struct {
	ChatActionId   int64  `json:"chatActionId"`
	ChatId         int64  `json:"chatId"`
	MessageType    string `json:"messageType"`
	CreateAt       string `json:"createAt"`
	ChatToken      string `json:"chatToken"`
	ChatMemberId   int64  `json:"chatMemberId"`
	MemberUUID     string `json:"memberUUID"`
	AccountType    string `json:"accountType"`
	AccountName    string `json:"accountName"`
	AccountUUID    string `json:"accountUUID"`
	AccountId      int64  `json:"accountId"`
	Open           bool   `json:"open"`
	MessageToken   string `json:"messageToken"`
	DisplayMessage string `json:"displayMessage"`
}

type Message struct {
	MessageID int64       `json:"messageId"`
	Type      string      `json:"type"`
	CreateAt  string      `json:"createAt"`
	Source    SourceData  `json:"source"`
	Message   MessageData `json:"message"`
}

type MessageData struct {
	MessageToken string `json:"messageToken"`
	MessageType  string `json:"messageType"`
	Message      string `json:"message"`
}

type SourceData struct {
	ChatMemberId int64  `json:"chatMemberId"`
	MemberUUID   string `json:"memberUUID"`
	AccountUUID  string `json:"accountUUID"`
	AccountId    int64  `json:"accountId"`
	ChatToken    string `json:"chatToken"`
	AccountName  string `json:"accountName"`
	AccountType  string `json:"accountType"`
}

type SendMessagePost struct {
	ChatId  int    `json:chatId`
	Type    string `json:type`
	Message string `json:message`
}

type MemberData struct {
	ChatMemberId int64  `json:"chatMemberId"`
	MemberUUID   string `json:"memberUUID"`
}

type MemberDataV2 struct {
	ChatMemberId int64  `json:"chatMemberId"`
	MemberUUID   string `json:"memberUUID"`
	AccountType  string `json:"accountType"`
}

type SendMessageResponse struct {
	Success      bool   `json:success`
	MessageId    int    `json:messageId`
	MessageToken string `json:messageToken`
}

type user struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

type ChatInfo struct {
	OrderId      int          `json:"orderId"`
	OrderToken   string       `json:"orderToken"`
	Status       string       `json:"status"`
	Note         string       `json:"note"`
	ChatId       int          `json:"chatId"`
	ChatToken    string       `json:"chatToken"`
	CreateDate   string       `json:"createDate"`
	LastTalkDate string       `json:"lastTalkDate"`
	Customer     CustomerData `json:"customer"`
}

type CustomerData struct {
	AccountId       int    `json:"accountId"`
	AccountUUID     string `json:"accountUUID"`
	Username        string `json:"username"`
	ProfileImageURL string `json:"profileImageURL"`
}
