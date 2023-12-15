package MerchantV1Main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
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
	err = db.QueryRow(query, tokenString, "seller").Scan(
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
	c.Set("sellerId", accessToken.AccountID)
	c.Next()

}

func extractTokenFromHeader(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

func CheckMerchantValid(c *gin.Context) {
	MerchantUUID := c.Param("MerchantUUID")
	if MerchantUUID == "" {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	sellerId, exists := c.Get("sellerId")
	if !exists {
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

	var merchant MerchantMerchantStruct.MerchantDetail
	query := `SELECT m.MerchantId, m.MerchantUUID, m.Name, m.Description,
				CASE WHEN m.Status = 'open' THEN true ELSE false END AS Open,
				CASE WHEN m.Visible = 'visible' THEN true ELSE false END AS Visible,
				m.ImageUrl, m.Address, m.Street, m.Building, 
				d.name_th as district, a.name_th as amphure, p.name_th as province, d.zip_code as zipcode, 
				m.CreateDate, m.UpdateDate
				FROM Merchant m
				LEFT JOIN districts d ON m.distric = d.id
				LEFT JOIN amphures a ON d.amphure_id = a.id
				LEFT JOIN provinces p ON a.province_id = p.id
				WHERE m.MerchantUUID = ? AND m.OwnerSellerId = ?`
	err = db.QueryRow(query, MerchantUUID, sellerId).Scan(
		&merchant.MerchantId,
		&merchant.MerchantUUID,
		&merchant.Name,
		&merchant.Description,
		&merchant.Open,
		&merchant.Visible,
		&merchant.ImageUrl,
		&merchant.Address.Address,
		&merchant.Address.Street,
		&merchant.Address.Building,
		&merchant.Address.District,
		&merchant.Address.Amphure,
		&merchant.Address.Province,
		&merchant.Address.ZipCode,
		&merchant.CreateDate,
		&merchant.UpdateDate,
	)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	merchant.ImageUrl = config.GetMainAPIURL() + merchant.ImageUrl

	c.Set("merchantData", merchant)
	c.Next()

}
