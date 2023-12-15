package MerchantV1Product

import (
	"net/http"

	"github.com/gin-gonic/gin"

	APIStructureMain "gobasic/struct"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
	UtilMerchant_Product "gobasic/util/merchant/product"
)

func GetProductInfo(c *gin.Context) {
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
	product, err := UtilMerchant_Product.GetProductDetail(ProductToken, mcData.MerchantId)
	if err != nil {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "ข้อมูลที่คุณระบุไม่ถูกต้อง โปรดตรวจสอบข้อมูลแล้วทำรายการใหม่อีกครั้ง")
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully get transaction",
		Data:       product,
	}
	c.JSON(http.StatusOK, respond)
}
