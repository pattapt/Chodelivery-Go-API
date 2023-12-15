package MerchantV1Product

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
	MerchantProductStruct "gobasic/struct/merchant/product"
	UtilMerchant_Product "gobasic/util/merchant/product"
)

func UploadProductImage(c *gin.Context) {
	merchantData, exists := c.Get("merchantData")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	mcData, ok := merchantData.(MerchantMerchantStruct.MerchantDetail)
	if !ok {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	ProductToken := c.Param("ProductToken")
	if ProductToken == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "ข้อมูลที่คุณระบุไม่ถูกต้อง โปรดตรวจสอบข้อมูลแล้วทำรายการใหม่อีกครั้ง")
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	// Product Info
	pddata, err := UtilMerchant_Product.GetProductDetail(ProductToken, mcData.MerchantId)
	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "ข้อมูลที่คุณระบุไม่ถูกต้อง โปรดตรวจสอบข้อมูลแล้วทำรายการใหม่อีกครั้ง")
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
	filePath := "cdn/products/" + newFileName
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

	// UPDATE PRODUCT
	filePath = "/" + filePath
	sql := `UPDATE Product SET ImageUrl = ?, UpdateDate = ?
			WHERE ProductId = ? AND MerchantId = ?`
	_, err = db.Exec(sql, filePath, config.GetCurrentDateTime(), pddata.ProductId, mcData.MerchantId)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully uploadImage of product",
		Data:       MerchantProductStruct.UpdateProductResponse{Success: true},
	}
	c.JSON(http.StatusOK, respond)
}
