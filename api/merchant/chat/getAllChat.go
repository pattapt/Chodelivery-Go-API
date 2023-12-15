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
)

func GetAllChat(c *gin.Context) {
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

	query := `SELECT t.OrderId, t.OrderToken, t.Status,
				t.Note, t.ChatId, c.ChatToken, t.CreateDate, c.LastTalkDate,
				cs.AccountId, cs.AccountUUID, cs.Username, cs.ProfileImageURL
				FROM Transaction t
				LEFT JOIN Merchant m ON t.MerchantId = m.MerchantId
				LEFT JOIN Account cs ON cs.AccountId = t.AccountId
				LEFT JOIN Chat c ON c.ChatId = t.ChatId
				WHERE t.MerchantId = ?
				ORDER BY c.LastTalkDate DESC
				LIMIT ?, ?;
	`
	rows, err := db.Query(query, merchantAccount.MerchantId, Start, limit)
	if err != nil {
		fmt.Print(err)
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer rows.Close()

	var ChatList []MerchantChatStruct.ChatInfo

	for rows.Next() {
		var chat MerchantChatStruct.ChatInfo
		err := rows.Scan(
			&chat.OrderId,
			&chat.OrderToken,
			&chat.Status,
			&chat.Note,
			&chat.ChatId,
			&chat.ChatToken,
			&chat.CreateDate,
			&chat.LastTalkDate,
			&chat.Customer.AccountId,
			&chat.Customer.AccountUUID,
			&chat.Customer.Username,
			&chat.Customer.ProfileImageURL,
		)
		if err != nil {
			fmt.Print(err)
			errorResponse := APIStructureMain.GeneralErrorMSG()
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		chat.Customer.ProfileImageURL = config.GetMainAPIURL() + chat.Customer.ProfileImageURL

		ChatList = append(ChatList, chat)
	}

	if len(ChatList) == 0 {
		// Return an empty array as the response
		respond := APIStructureMain.Response{
			StatusCode: http.StatusOK,
			Status:     "S",
			Message:    "No chat found",
			Data:       []MerchantChatStruct.Message{},
		}
		c.JSON(http.StatusOK, respond)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully retrieved chat",
		Data:       ChatList,
	}
	c.JSON(http.StatusOK, respond)
}
