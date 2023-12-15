package MerchantV1Store

import (
	"net/http"

	"github.com/gin-gonic/gin"

	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
)

func UpdateMerchant(c *gin.Context) {
	MerchantUUID := c.Param("MerchantUUID")
	if MerchantUUID == "" {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer db.Close()

	_, exists := c.Get("sellerId")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var postData MerchantMerchantStruct.EditMerchantProfilePost

	if err := c.ShouldBindJSON(&postData); err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// CHECK IS District EXIST
	var district int
	err = db.QueryRow(`SELECT id FROM districts WHERE id = ?`, postData.District).Scan(&district)
	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลที่อยู่ไม่ถูกต้อง", "กรุณาตรวจสอบข้อมูลที่อยู่อีกครั้ง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Name == "" || postData.Description == "" || postData.Address == "" || postData.Street == "" || postData.Building == "" || postData.District == 0 {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ครบถ้วน", "ข้อมูลที่ท่านระบุไม่ครบถ้วน โปรดตรวจสอบข้อมูลแล้วทำรายการใหม่อีกครั้ง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.PromptpayPhone == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุหมายเลขพร้อมเพย์", "กรุณาระบุหมายเลขพร้อมเพย์สำหรับใช้ในการรับเงิน")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// UPDATE MERCHANT
	Status := "open"
	if postData.Open {
		Status = "open"
	} else {
		Status = "close"
	}
	Visible := "visible"
	if postData.Visible {
		Visible = "visible"
	} else {
		Visible = "invisible"
	}

	sql := `UPDATE Merchant SET Name = ?, Description = ?, PromptpayPhone = ?, Address = ?,
				Street = ?, Building = ?, distric = ?, Status = ?, Visible = ? WHERE MerchantUUID = ?`
	_, err = db.Exec(sql, postData.Name, postData.Description, postData.PromptpayPhone, postData.Address,
		postData.Street, postData.Building, postData.District, Status, Visible, MerchantUUID)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Get Profile
	respond := APIStructureMain.Response{
		StatusCode: 200,
		Status:     "S",
		Message:    "Successfully Update Data",
		Data:       MerchantMerchantStruct.UpdateProfileResponse{Success: true},
	}
	c.JSON(http.StatusOK, respond)

}
