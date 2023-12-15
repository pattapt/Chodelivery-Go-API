package MerchantV1Chat

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	UtilStore_Chat "gobasic/util/store/chat"

	MerchantChatStruct "gobasic/struct/merchant/chat"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
	// MerchantChatStruct "gobasic/struct/merchant/chat"
)

func ViewImage(c *gin.Context) {
	merchantData, exists := c.Get("merchantData")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data 1")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	_, ok := merchantData.(MerchantMerchantStruct.MerchantDetail)
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

	_, err = UtilStore_Chat.GetMemberAccount(sellerId.(int), chat.ChatId, "seller")
	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("คุณไม่สามารถทำรายการได้", "คุณไม่สามารถดูรูปภาพนี้ได้")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	ImageToken := c.Param("ImageToken")
	if ImageToken == "" {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	db, err = database.ConnectDB()
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer db.Close()

	// GET FILENAME BY TOKEN
	FileName := ""
	sql := `SELECT ImageURL AS FileName FROM ChatImage WHERE Token = ?`
	err = db.QueryRow(sql, ImageToken).Scan(&FileName)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Construct the file path
	filePath := filepath.Join(FileName)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer file.Close()

	// Determine content type based on file extension
	contentType := getContentType(FileName)

	// Read the file data
	fileData, err := io.ReadAll(file)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Respond with the file data
	c.Data(http.StatusOK, contentType, fileData)
}

// getContentType determines the content type based on the file extension
func getContentType(filename string) string {
	extension := strings.ToLower(filepath.Ext(filename))
	switch extension {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		// Default to "application/octet-stream" for unknown file types
		return "application/octet-stream"
	}
}
