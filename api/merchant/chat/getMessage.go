package MerchantV1Chat

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"

	MerchantChatStruct "gobasic/struct/merchant/chat"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
	// MerchantChatStruct "gobasic/struct/merchant/chat"
)

func GetMessage(c *gin.Context) {
	merchantData, exists := c.Get("merchantData")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data 1")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	merchantAccount, ok := merchantData.(MerchantMerchantStruct.MerchantDetail)
	if !ok {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data 2")
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
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data 3")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	ChatData, exists := c.Get("chatData")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data 4")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	chat, ok := ChatData.(MerchantChatStruct.ChatData)
	if !ok {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data 5")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	pageStr := c.DefaultQuery("page", "1")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 50
	}

	Start := (page - 1) * limit

	query := `
		SELECT ca.ChatActionId AS MessageID, c.ChatToken, ca.MessageType, ca.CreateAt,
		cm.MemberUUID, cm.ChatMemberId, cm.AccountType,
		CONCAT_WS('-', ac.Username, CONCAT(sl.Name, " ", sl.LastName)) AS AccountName,
		CONCAT_WS('-', ac.AccountUUID, sl.SellerUUID) AS AccountUUID,
		CONCAT_WS('-', ac.AccountId, sl.SellerId) AS AccountId,
		CASE
			WHEN ca.MessageType = 'Message' THEN ctm.Token
			WHEN ca.MessageType = 'Image' THEN cti.Token
		END AS MessageToken,
		CASE
			WHEN ca.MessageType = 'Message' THEN ctm.Text
			WHEN ca.MessageType = 'Image' THEN cti.Token
		END AS Message,
		CASE
			WHEN cm.AccountType = ? AND sl.SellerId = ? THEN "messageSend" ELSE "messageReceive"
		END AS type
		FROM ChatAction ca
		LEFT JOIN Chat c ON c.ChatId = ca.ChatId
		LEFT JOIN ChatMember cm ON cm.ChatMemberId = ca.MemberId
		LEFT JOIN Account ac ON ac.AccountId = cm.AccountId AND cm.AccountType = 'customer'
		LEFT JOIN Seller sl ON sl.SellerId = cm.AccountId AND cm.AccountType = 'seller'
		LEFT JOIN ChatMessage ctm ON ctm.MessageId = ca.RefId AND ca.MessageType = 'Message'
		LEFT JOIN ChatImage cti ON cti.ImageId = ca.RefId AND ca.MessageType = 'Image'
		WHERE ca.ChatId = ?
		ORDER BY ca.ChatActionId DESC
		LIMIT ?, ?
	`
	rows, err := db.Query(query, "seller", sellerId, chat.ChatId, Start, limit)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer rows.Close()

	var Messages []MerchantChatStruct.Message

	for rows.Next() {
		var Message MerchantChatStruct.Message
		err := rows.Scan(
			&Message.MessageID,
			&Message.Source.ChatToken,
			&Message.Message.MessageType,
			&Message.CreateAt,
			&Message.Source.MemberUUID,
			&Message.Source.ChatMemberId,
			&Message.Source.AccountType,
			&Message.Source.AccountName,
			&Message.Source.AccountUUID,
			&Message.Source.AccountId,
			&Message.Message.MessageToken,
			&Message.Message.Message,
			&Message.Type,
		)
		if err != nil {
			fmt.Print(err)
			errorResponse := APIStructureMain.GeneralErrorMSG()
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}

		if Message.Message.MessageType == "Image" {
			Message.Message.Message = config.GetMainAPIURL() + "/api/merchant/v1/store/" + merchantAccount.MerchantUUID + "/chat/" + chat.ChatToken + "/Image/" + Message.Message.Message
		}

		Messages = append(Messages, Message)
	}

	if len(Messages) == 0 {
		// Return an empty array as the response
		respond := APIStructureMain.Response{
			StatusCode: http.StatusOK,
			Status:     "S",
			Message:    "No message found for the chat",
			Data:       []MerchantChatStruct.Message{},
		}
		c.JSON(http.StatusOK, respond)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully retrieved message",
		Data:       Messages,
	}
	c.JSON(http.StatusOK, respond)
}
