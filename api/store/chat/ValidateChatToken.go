package StoreV1Chat

import (
	"net/http"

	"github.com/gin-gonic/gin"

	database "gobasic/database"
	APIStructureMain "gobasic/struct"

	UtilStore_Chat "gobasic/util/store/chat"
)

func ValidateChat(c *gin.Context) {
	ChatToken := c.Param("ChatToken")
	if ChatToken == "" {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	// Check AccessToken table
	db, err := database.ConnectDB()
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer db.Close()

	UtilStore_Chat, err := UtilStore_Chat.GetChatInfoByToken(ChatToken)
	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("ไม่พบข้อมูล", "ระบบไม่พบข้อมูลแชทที่คุณต้องการ กรุณาตรวจสอบข้อมูลแล้วลองใหม่อีกครั้ง")
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	// Check if the token is expired
	if UtilStore_Chat.ChatId == 0 {
		errorResponse := APIStructureMain.ErrorRespondAuth("ไม่พบข้อมูล", "ระบบไม่พบข้อมูลแชทที่คุณต้องการ กรุณาตรวจสอบข้อมูลแล้วลองใหม่อีกครั้ง")
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	c.Set("chatData", UtilStore_Chat)
	c.Next()

}
