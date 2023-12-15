package MerchantV1Auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantAccountStruct "gobasic/struct/merchant/account"
)

func LoginAccount(c *gin.Context) {
	var postData MerchantAccountStruct.LoginAccountPost

	if err := c.ShouldBindJSON(&postData); err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Email == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุ Email", "กรุณาระบุ Email ของท่านเพื่อทำการเข้าสู่ระบบ")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Password == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุรหัสผ่าน", "กรุณาระบุรหัสผ่านของท่านเพื่อทำการเข้าสู่ระบบ")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if config.IsValidEmail(postData.Email) == false {
		errorResponse := APIStructureMain.ErrorRespondAuth("อีเมล์ไม่ถูกต้อง", "รูปแบบอีเมล์ที่คุณระบุไม่ถูกต้อง โปรดตรวจสอบอีเมล์ให้ถูกต้องแล้วทำรายการใหม่อีกครั้ง")
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

	var account MerchantAccountStruct.SellerAccount
	query := "SELECT * FROM Seller WHERE Email = ?"
	err = db.QueryRow(query, postData.Email).Scan(&account.SellerID, &account.SellerUUID, &account.Role, &account.Name, &account.LastName, &account.Email, &account.Password, &account.MerchantID, &account.ProfileImageURL, &account.RegisterDate, &account.LastLoginDate, &account.IP)
	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "อีเมล์หรือรหัสผ่านที่คุณระบุไม่ถูกต้อง โปรดตรวจสอบแล้วลองใหม่ภายหลัง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if err == sql.ErrNoRows {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "อีเมล์หรือรหัสผ่านที่คุณระบุไม่ถูกต้อง โปรดตรวจสอบแล้วลองใหม่ภายหลัง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(postData.Password))
	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "อีเมล์หรือรหัสผ่านที่คุณระบุไม่ถูกต้อง โปรดตรวจสอบแล้วลองใหม่ภายหลัง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	RefreshToken := "R-" + config.GenerateRefreshToken()
	CreateDate := config.GetCurrentDateTime()
	currentDate := time.Now()
	expiredDate := currentDate.AddDate(1, 0, 0)
	expiredDateFormatted := expiredDate.Format("2006-01-02 15:04:05")

	sql := `INSERT INTO AccessToken (AccountId, TokenType, AccountType, Token, IssueDate, ExpiredDate, Ip)
			VALUES (?, ?, ?, ?, ?, ?, ?)`
	Result, err := db.Exec(sql, account.SellerID, "refreshToken", "seller", RefreshToken, CreateDate, expiredDateFormatted, c.ClientIP())
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	RefreshTokenID, err := Result.LastInsertId()
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	AccessToken := "A-" + config.GenerateRefreshToken()
	expiredDateAccessToken := currentDate.AddDate(0, 0, 1)
	expiredDateAccTokenFormatted := expiredDateAccessToken.Format("2006-01-02 15:04:05")
	sql = `INSERT INTO AccessToken (AccountId, TokenType, AccountType, RefreshTokenId, Token, IssueDate, ExpiredDate, Ip)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	Result, err = db.Exec(sql, account.SellerID, "accessToken", "seller", RefreshTokenID, AccessToken, CreateDate, expiredDateAccTokenFormatted, c.ClientIP())
	if err != nil {
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
			SellerID:        account.SellerID,
			SellerUUID:      account.SellerUUID,
			Role:            account.Role,
			Name:            account.Name,
			LastName:        account.LastName,
			Email:           account.Email,
			ProfileImageURL: config.GetMainAPIURL() + account.ProfileImageURL,
			RefreshToken:    RefreshToken,
			AccessToken:     AccessToken,
			CreateDate:      CreateDate,
			ExpiredDate:     expiredDateAccTokenFormatted,
		},
	}
	c.JSON(http.StatusOK, respond)

}
