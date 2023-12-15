package StoreV1Auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	StoreAccountStruct "gobasic/struct/store/account"
)

func LoginAccount(c *gin.Context) {
	var postData StoreAccountStruct.LoginAccountPost

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

	var account StoreAccountStruct.AccountData
	query := "SELECT * FROM Account WHERE Email = ?"
	err = db.QueryRow(query, postData.Email).Scan(&account.AccountID, &account.AccountUUID, &account.Username, &account.Email, &account.Password, &account.ProfileImageURL, &account.RegisterDate, &account.LastLoginDate)
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
	Result, err := db.Exec(sql, account.AccountID, "refreshToken", "customer", RefreshToken, CreateDate, expiredDateFormatted, c.ClientIP())
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
	Result, err = db.Exec(sql, account.AccountID, "accessToken", "customer", RefreshTokenID, AccessToken, CreateDate, expiredDateAccTokenFormatted, c.ClientIP())
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
		Data: APIStructureMain.AccountDataLogin{
			AccountID:   account.AccountID,
			AccountUUID: account.AccountUUID,
			Username:    account.Username,
			Email:       account.Email,
			ProfileImageURL: func() string {
				if account.ProfileImageURL.Valid {
					return config.GetMainAPIURL() + account.ProfileImageURL.String
				}
				return ""
			}(),
			RefreshToken: RefreshToken,
			AccessToken:  AccessToken,
			CreateDate:   CreateDate,
			ExpiredDate:  expiredDateAccTokenFormatted,
		},
	}
	c.JSON(http.StatusOK, respond)

}
