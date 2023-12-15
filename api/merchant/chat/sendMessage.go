package MerchantV1Chat

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pusher/pusher-http-go/v5"

	"gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	UtilStore_Chat "gobasic/util/store/chat"

	MerchantChatStruct "gobasic/struct/merchant/chat"

	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
)

func SendMessage(c *gin.Context) {
	merchantData, exists := c.Get("merchantData")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	_, ok := merchantData.(MerchantMerchantStruct.MerchantDetail)
	if !ok {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer db.Close()

	sellerId, exists := c.Get("sellerId")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	ChatData, exists := c.Get("chatData")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var postData MerchantChatStruct.SendMessagePost

	if err := c.ShouldBindJSON(&postData); err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Type != "Message" {
		errorResponse := APIStructureMain.ErrorRespondAuth("ชนิดของข้อมูลไม้ถูกต้อง", "ข้อมูลที่คุณระบุไม่ถูกต้อง โปรดตรวจสอบข้อมูลแล้วดำเนินการใหม่อีกครั้ง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Message == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุข้อความ", "กรุณาระบุข้อความที่ต้องการส่ง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	chat, ok := ChatData.(MerchantChatStruct.ChatData)
	if !ok {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	Member, err := UtilStore_Chat.GetMemberAccount(sellerId.(int), chat.ChatId, "seller")
	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("คุณไม่สามารถทำรายการได้", "คุณไม่ได้รับอนุญาติให้ส่งข่้อความในห้องแชทนี้")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	MessageToken := config.GenerateRefreshToken()
	DateTime := config.GetCurrentDateTime()
	// ADD MESSAGE TO CHAT MESSAGE
	sql := `INSERT INTO ChatMessage (Token, Text, CreateAt) VALUE(?, ?, ?)`
	Result, err := db.Exec(sql, MessageToken, postData.Message, DateTime)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	MessageId, err := Result.LastInsertId()

	// CREATE ACTION
	sql = `INSERT INTO ChatAction (ChatId, MemberId, RefId, MessageType, CreateAt) VALUE(?, ?, ?, ?, ?)`
	Result, err = db.Exec(sql, chat.ChatId, Member.ChatMemberId, MessageId, "Message", DateTime)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	ActionId, err := Result.LastInsertId()

	// UPDATE CHAT LAST TALK
	sql = `UPDATE Chat SET LastTalkDate = ? WHERE ChatId = ?`
	_, err = db.Exec(sql, DateTime, chat.ChatId)

	// BOARDCAST TO PUSHER
	var client = pusher.Client{
		AppID:   "1721916",
		Key:     "8dbaa73ce7ec51885a7e",
		Secret:  "3b1f5bd8a734605691bd",
		Cluster: "ap1",
		Secure:  true,
	}

	MembersList, err := UtilStore_Chat.GetMemberInChat(chat.ChatId)
	// fmt.Print(MembersList)
	for _, Memberx := range MembersList {
		Data := MerchantChatStruct.Message{
			MessageID: int64(ActionId),
			Type:      "messageReceive",
			CreateAt:  DateTime,
			Source: MerchantChatStruct.SourceData{
				ChatToken:   chat.ChatToken,
				MemberUUID:  Member.MemberUUID,
				AccountUUID: Member.AccountUUID,
				AccountId:   Member.AccountId,
				AccountName: Member.AccountName,
				AccountType: Member.AccountType,
			},
			Message: MerchantChatStruct.MessageData{
				MessageToken: MessageToken,
				MessageType:  "Message",
				Message:      postData.Message,
			},
		}
		if Member.MemberUUID == Memberx.MemberUUID {
			Data.Type = "messageSend"
		}
		client.Trigger(chat.ChatToken, Memberx.MemberUUID, Data)
		jsonData, err := json.Marshal(Data)
		if err != nil {
			return
		}
		UtilStore_Chat.BoardcastSocket(chat.ChatToken+":"+Memberx.MemberUUID, string(jsonData))
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully save message",
		Data: MerchantChatStruct.SendMessageResponse{
			Success:      true,
			MessageId:    int(ActionId),
			MessageToken: MessageToken,
		},
	}
	c.JSON(http.StatusOK, respond)
}
