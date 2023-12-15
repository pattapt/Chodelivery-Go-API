package MerchantV1Store

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
)

func UpdateMerchantImage(c *gin.Context) {
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

	err = c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	file, fileHeader, err := c.Request.FormFile("Image")
	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาเลือกรูปภาพ", "คุณไม่ได้ทำการเลือกรูปภาพ กรุณาทำการเลือกรูปภาพแล้วดำเนินการใหม่อีกครั้ง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer file.Close()

	// Handle the file (e.g., save to disk)
	fileName := fileHeader.Filename
	fileExt := filepath.Ext(fileName)
	newFileName := config.GenerateRefreshToken() + fileExt
	filePath := "cdn/profile/merchant/" + newFileName
	out, err := os.Create(filePath)
	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("เกิดข้อผิดพลาด", "ระบบเกิดข้อผิดพลาดในการอัพโหลดรูปภาพ โปรดทำรายการใหม่ภายหลัง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("เกิดข้อผิดพลาด", "ระบบเกิดข้อผิดพลาดในการอัพโหลดรูปภาพ โปรดทำรายการใหม่ภายหลัง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	filePath = "/" + filePath
	sql := `UPDATE Merchant SET ImageUrl = ? WHERE MerchantUUID = ?`
	_, err = db.Exec(sql, filePath, MerchantUUID)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Get Profile
	respond := APIStructureMain.Response{
		StatusCode: 200,
		Status:     "S",
		Message:    "Successfully Update Image",
		Data:       MerchantMerchantStruct.UpdateProfileResponse{Success: true},
	}
	c.JSON(http.StatusOK, respond)

}
