package MerchantV1Auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	struct_api "gobasic/struct"
	MerchantAccountStruct "gobasic/struct/merchant/account"
)

func RegisterAccount(c *gin.Context) {
	var postData MerchantAccountStruct.RegisterAccountPost

	if err := c.ShouldBindJSON(&postData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// fmt.Printf("Email: %s\n", postData.Email)
	// fmt.Printf("Password: %s\n", postData.Password)

	if postData.Name == "" {
		errorResponse := struct_api.ErrorRespondAuth("กรุณาระบุชื่อของคุณ", "กรุณาระบุชื่อของท่านเพื่อทำการสมัครสมาชิก")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.LastName == "" {
		errorResponse := struct_api.ErrorRespondAuth("กรุณาระบุนามสกุล", "กรุณาระบุนามสกุลของท่านเพื่อทำการสมัครสมาชิก")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Email == "" {
		errorResponse := struct_api.ErrorRespondAuth("กรุณาระบุ Email", "กรุณาระบุ Email ของท่านเพื่อทำการสมัครสมาชิก")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Password == "" {
		errorResponse := struct_api.ErrorRespondAuth("กรุณาระบุรหัสผ่าน", "กรุณาระบุรหัสผ่านของท่านเพื่อทำการสมัครสมาชิก")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if config.IsValidEmail(postData.Email) == false {
		errorResponse := struct_api.ErrorRespondAuth("อีเมล์ไม่ถูกต้อง", "รูปแบบอีเมล์ที่คุณระบุไม่ถูกต้อง โปรดตรวจสอบอีเมล์ให้ถูกต้องแล้วทำรายการใหม่อีกครั้ง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		errorResponse := struct_api.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer db.Close()

	var emailExists bool
	query := "SELECT COUNT(*) FROM Seller WHERE Email = ?"
	err = db.QueryRow(query, postData.Email).Scan(&emailExists)
	if err != nil {
		errorResponse := struct_api.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if emailExists {
		errorResponse := struct_api.ErrorRespondAuth("อีเมล์นี้ถูกใช้งานแล้ว", "อีเมล์ที่ท่านระบุถูกใช้งานโดยผู้ใช้งานท่านอื่นแล้ว โปรดใช้อีเมล์อื่นในการสมัครสมาชิก")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	accountUUID := uuid.New().String()
	// Decrypt and hash the password (assuming it's a SHA-256 hash)
	// You should replace this with your actual password decryption logic
	decryptedPassword := postData.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(decryptedPassword), bcrypt.DefaultCost)
	if err != nil {
		errorResponse := struct_api.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Format the RegisterDate and LastLoginDate
	registerDate := config.GetCurrentDateTime()
	lastLoginDate := config.GetCurrentDateTime()

	// Insert the data into the Account table
	insertQuery := `
        INSERT INTO Seller (SellerUUID, Role, Name, LastName, Email, password, ProfileImageURL, RegisterDate, LastLoginDate, ip)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	_, err = db.Exec(insertQuery, accountUUID, postData.Role, postData.Name, postData.LastName, postData.Email, hashedPassword, "/cdn/profile/merchant/default.png", registerDate, lastLoginDate, c.ClientIP())
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: 200,
		Status:     "S",
		Message:    "Successfully Register",
		Data:       APIStructureMain.RegisterAccountStoreSuccess{UUID: accountUUID, Username: postData.Name + " " + postData.LastName, Email: postData.Email},
	}
	c.JSON(http.StatusOK, respond)
}

func RegisterToType(c *gin.Context) {

}
