package MerchantV1Chat

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pusher/pusher-http-go/v5"

	"gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	UtilStore_Chat "gobasic/util/store/chat"

	MerchantChatStruct "gobasic/struct/merchant/chat"

	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
)

func SendImage(c *gin.Context) {
	merchantData, exists := c.Get("merchantData")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	merchantAccount, ok := merchantData.(MerchantMerchantStruct.MerchantDetail)
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

	err = c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	file, fileHeader, err := c.Request.FormFile("Image")
	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาเลือกรูปภาพ", "คุณไม่ได้ทำการเลือกรูปภาพ กรุณาทำการเลือกรูปภาพแล้วดำเนินการใหม่อีกครั้ง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer file.Close()

	// Handle the file (e.g., save to disk)
	fileName := fileHeader.Filename
	fileExt := filepath.Ext(fileName)
	newFileName := config.GenerateRefreshToken() + fileExt
	filePath := "cdn/chat/" + newFileName
	out, err := os.Create(filePath)
	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("เกิดข้อผิดพลาด", "ระบบเกิดข้อผิดพลาดในการอัพโหลดรูปภาพ โปรดทำรายการใหม่ภายหลัง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("เกิดข้อผิดพลาด", "ระบบเกิดข้อผิดพลาดในการอัพโหลดรูปภาพ โปรดทำรายการใหม่ภายหลัง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	MessageToken := config.GenerateRefreshToken()
	DateTime := config.GetCurrentDateTime()
	// ADD MESSAGE TO CHAT MESSAGE
	sql := `INSERT INTO ChatImage (Token, ImageURL, CreateAt) VALUE(?, ?, ?)`
	Result, err := db.Exec(sql, MessageToken, filePath, DateTime)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	ImageId, err := Result.LastInsertId()

	// CREATE ACTION
	sql = `INSERT INTO ChatAction (ChatId, MemberId, RefId, MessageType, CreateAt) VALUE(?, ?, ?, ?, ?)`
	Result, err = db.Exec(sql, chat.ChatId, Member.ChatMemberId, ImageId, "Image", DateTime)
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
				Message:      config.GetMainAPIURL() + "/api/merchant/v1/store/" + merchantAccount.MerchantUUID + "/chat/" + chat.ChatToken + "/Image/" + MessageToken,
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
