package StoreV1Main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	StoreAccountStruct "gobasic/struct/store/account"
)

func AccessTokenMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Authorization header is missing")
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	tokenString := extractTokenFromHeader(authHeader)
	if tokenString == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Authorization header is missing")
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

	var accessToken StoreAccountStruct.AccessTokenData
	query := "SELECT * FROM AccessToken WHERE TokenType = 'accessToken' AND Token = ? AND AccountType = ?"
	err = db.QueryRow(query, tokenString, "customer").Scan(
		&accessToken.AccessTokenID,
		&accessToken.AccountID,
		&accessToken.TokenType,
		&accessToken.AccountType,
		&accessToken.RefreshTokenID,
		&accessToken.Token,
		&accessToken.IssueDate,
		&accessToken.ExpiredDate,
		&accessToken.IP,
	)

	if err != nil {
		errorResponse := APIStructureMain.AccessTokenNotValid()
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	// Check if the token is expired
	expiredDate, err := time.Parse("2006-01-02 15:04:05", accessToken.ExpiredDate)
	if err != nil || time.Now().After(expiredDate) {
		errorResponse := APIStructureMain.AccessTokenNotValid()
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	// Set user information in the context for further use in handlers
	c.Set("userID", accessToken.AccountID)
	c.Next()

}

func extractTokenFromHeader(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

func GetMerchantId(c *gin.Context) {
	MerchantToken := c.Param("merchantToken")
	if MerchantToken == "" {
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

	MerchantId := 0
	query := "SELECT MerchantId FROM Merchant WHERE MerchantUUID = ?"
	err = db.QueryRow(query, MerchantToken).Scan(&MerchantId)

	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("ไม่พบข้อมูล", "ระบบไม่พบข้อมูลร้านค้าที่คุณต้องการ กรุณาตรวจสอบข้อมูลแล้วลองใหม่อีกครั้ง")
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	// Check if the token is expired
	if MerchantId == 0 {
		errorResponse := APIStructureMain.ErrorRespondAuth("ไม่พบข้อมูล", "ระบบไม่พบข้อมูลร้านค้าที่คุณต้องการ กรุณาตรวจสอบข้อมูลแล้วลองใหม่อีกครั้ง")
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	// Set user information in the context for further use in handlers
	c.Set("merchantId", MerchantId)
	c.Next()

}
