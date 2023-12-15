package MerchantV1Auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantAccountStruct "gobasic/struct/merchant/account"
	StoreAccountStruct "gobasic/struct/store/account"
)

func InvokeAccessToken(c *gin.Context) {
	var postData StoreAccountStruct.InvokeAccessToken

	if err := c.ShouldBindJSON(&postData); err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.RefreshToken == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.UUID == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth()
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

	var RefreshTokenData MerchantAccountStruct.RefreshTokenData
	query := `
		SELECT at.AccessTokenId, sl.SellerId AS accID, at.Token, at.TokenType, at.ExpiredDate,
		sl.SellerUUID, sl.Name, sl.LastName, sl.Email, sl.Password, sl.ProfileImageUrl
		FROM AccessToken at
		LEFT JOIN Seller sl ON at.AccountId = sl.SellerId
		WHERE at.Token = ? AND at.TokenType = ? AND at.AccountType = ?
	`
	err = db.QueryRow(query, postData.RefreshToken, "refreshToken", "seller").Scan(
		&RefreshTokenData.RefreshTokenID,
		&RefreshTokenData.SellerID,
		&RefreshTokenData.Token,
		&RefreshTokenData.TokenType,
		&RefreshTokenData.ExpiredDate,
		&RefreshTokenData.SellerUUID,
		&RefreshTokenData.Name,
		&RefreshTokenData.LastName,
		&RefreshTokenData.Email,
		&RefreshTokenData.Password,
		&RefreshTokenData.ProfileImageURL)
	if err != nil {
		fmt.Println(err)
		errorResponse := APIStructureMain.ErrorRespondAuth()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if err == sql.ErrNoRows {
		errorResponse := APIStructureMain.ErrorRespondAuth()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	CreateDate := config.GetCurrentDateTime()
	currentDate := time.Now()
	AccessToken := "A-" + config.GenerateRefreshToken()
	expiredDateAccessToken := currentDate.AddDate(0, 0, 1)
	expiredDateAccTokenFormatted := expiredDateAccessToken.Format("2006-01-02 15:04:05")
	sql := `INSERT INTO AccessToken (AccountId, TokenType, AccountType, RefreshTokenId, Token, IssueDate, ExpiredDate, Ip)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(sql, RefreshTokenData.SellerID, "accessToken", "seller", RefreshTokenData.RefreshTokenID, AccessToken, CreateDate, expiredDateAccTokenFormatted, c.ClientIP())
	if err != nil {
		fmt.Println(err)
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Authentication successful
	respond := APIStructureMain.Response{
		StatusCode: 200,
		Status:     "S",
		Message:    "Successfully Register",
		Data: MerchantAccountStruct.SellerAccountDataLogin{
			SellerID:   RefreshTokenData.SellerID,
			SellerUUID: RefreshTokenData.SellerUUID,
			Name:       RefreshTokenData.Name,
			LastName:   RefreshTokenData.LastName,
			Email:      RefreshTokenData.Email,
			ProfileImageURL: func() string {
				if RefreshTokenData.ProfileImageURL.Valid {
					return config.GetMainAPIURL() + RefreshTokenData.ProfileImageURL.String
				}
				return ""
			}(),
			AccessToken: AccessToken,
			CreateDate:  CreateDate,
			ExpiredDate: expiredDateAccTokenFormatted,
		},
	}

	c.JSON(http.StatusOK, respond)

}
