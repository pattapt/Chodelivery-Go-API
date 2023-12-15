package MerchantV1Store

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantProfileStruct "gobasic/struct/merchant/merchant"
)

func Createmerchant(c *gin.Context) {

	db, err := database.ConnectDB()
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer db.Close()

	sellerId, exists := c.Get("sellerId")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var postData MerchantProfileStruct.CreateMerchantRequest

	if err := c.ShouldBindJSON(&postData); err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Name == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุชื่อร้านค้า", "กรุณาระบุชื่อร้านค้าของท่าน")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Description == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุคำอธิบายเกี่ยวกับร้านค้า", "กรุณาระบุคำอธิบายเกี่ยวกับร้านค้าของคุณสักเล็กน้อย")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Address == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุที่อยู่", "กรุณาระบุที่อยู่ของร้านค้า")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Street == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุถนน", "กรุณาระบุถนนที่ร้านค้าตั้งอยู่")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Building == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุอาคาร", "กรุณาระบุอาคารที่ร้านค้าตั้งอยู่")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.District == 0 {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุตำบล", "กรุณาระบุตำบล เพื่อดำเนินการต่อ")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Amphures == 0 {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุอำเภอ", "กรุณาระบุอำเภอ เพื่อดำเนินการต่อ")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Provinces == 0 {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุจังหวัด", "กรุณาระบุจังหวัด เพื่อดำเนินการต่อ")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.PromptpayId == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุหมายเลขพร้อมเพย์", "กรุณาระบุหมายเลขพร้อมเพย์สำหรับใช้ในการรับเงิน")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	total := 0
	query := "SELECT COUNT(1) AS total FROM Seller WHERE SellerId = ? AND MerchantId IS NULL"
	err = db.QueryRow(query, sellerId).Scan(&total)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if total == 0 {
		errorResponse := APIStructureMain.ErrorRespondAuth("ไม่สามารถดำเนินการได้", "ระบบไม่สามารถทำการสร้างร้านค้าใหม่ให้กับท่านในตอนนี้ เนื่องจากขณะนี้จำกัดการสร้างร้านไว้เพียง 1 ร้านต่อ 1 Account")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Create Merchant
	merchantUUID := uuid.New().String()
	DateTime := config.GetCurrentDateTime()
	insertQuery := `INSERT INTO Merchant (MerchantId, MerchantUUID, OwnerSellerId, Name, Description, Status, ImageUrl, Visible, PromptpayPhone, Address, Street, Building, distric, CreateDate, UpdateDate, CreateIp) 
			VALUES (NULL, ?, ?, ?, ?, 'open', ?, 'visible', ?, ?, ?, ?, ?, ?, ?, ?);`
	Result, err := db.Exec(insertQuery, merchantUUID, sellerId, postData.Name, postData.Description, "/cdn/profile/customer/default.png", postData.PromptpayId, postData.Address, postData.Street, postData.Building, postData.District, DateTime, DateTime, c.ClientIP())
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Update MerchantID to Seller
	MerchantId, err := Result.LastInsertId()
	updateQuery := `UPDATE Seller SET MerchantId = ? WHERE SellerId = ?`
	_, err = db.Exec(updateQuery, MerchantId, sellerId)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Get Profile
	respond := APIStructureMain.Response{
		StatusCode: 200,
		Status:     "S",
		Message:    "Successfully create merchant",
		Data: map[string]interface{}{
			"success":      true,
			"merchantUUID": merchantUUID,
		},
	}
	c.JSON(http.StatusOK, respond)

}
