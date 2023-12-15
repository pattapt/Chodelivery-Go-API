package StoreV1Auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	StoreAccountStruct "gobasic/struct/store/account"
)

func RegisterAccount(c *gin.Context) {
	var postData StoreAccountStruct.RegisterAccountPost

	if err := c.ShouldBindJSON(&postData); err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	// fmt.Printf("Email: %s\n", postData.Email)
	// fmt.Printf("Password: %s\n", postData.Password)
	if postData.Username == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุชื่อผู้ใข้งาน", "กรุณาระบุชื่อผู้ใข้งานของท่านเพื่อทำการสมัครสมาชิก")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Email == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุ Email", "กรุณาระบุ Email ของท่านเพื่อทำการสมัครสมาชิก")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Password == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุรหัสผ่าน", "กรุณาระบุรหัสผ่านของท่านเพื่อทำการสมัครสมาชิก")
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

	var emailExists bool
	query := "SELECT COUNT(*) FROM Account WHERE Email = ?"
	err = db.QueryRow(query, postData.Email).Scan(&emailExists)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if emailExists {
		errorResponse := APIStructureMain.ErrorRespondAuth("อีเมล์นี้ถูกใช้งานแล้ว", "อีเมล์ที่ท่านระบุถูกใช้งานโดยผู้ใช้งานท่านอื่นแล้ว โปรดใช้อีเมล์อื่นในการสมัครสมาชิก")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	accountUUID := uuid.New().String()
	// Decrypt and hash the password (assuming it's a SHA-256 hash)
	// You should replace this with your actual password decryption logic
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(postData.Password), bcrypt.DefaultCost)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Format the RegisterDate and LastLoginDate
	registerDate := config.GetCurrentDateTime()
	lastLoginDate := config.GetCurrentDateTime()

	// Insert the data into the Account table
	insertQuery := `
        INSERT INTO Account (AccountUUID, Username, Email, Password, RegisterDate, LastLoginDate)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err = db.Exec(insertQuery, accountUUID, postData.Username, postData.Email, hashedPassword, registerDate, lastLoginDate)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: 200,
		Status:     "S",
		Message:    "Successfully Register",
		Data:       APIStructureMain.RegisterAccountStoreSuccess{UUID: accountUUID, Username: postData.Username, Email: postData.Email},
	}
	c.JSON(http.StatusOK, respond)
}

func RegisterToType(c *gin.Context) {

}
