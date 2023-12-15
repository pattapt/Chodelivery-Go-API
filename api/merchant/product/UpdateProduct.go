package MerchantV1Product

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
	MerchantProductStruct "gobasic/struct/merchant/product"
	UtilMerchant_Product "gobasic/util/merchant/product"
)

func UpdateProduct(c *gin.Context) {
	var postData MerchantProductStruct.EditProductPost

	if err := c.ShouldBindJSON(&postData); err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.ProductId == 0 || postData.Name == "" || postData.Description == "" || postData.Price == 0 || postData.Cost == 0 || postData.Status == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("กรุณาระบุข้อมูลให้ครบถ้วน", "กรุณาระบข้อมูลให้ครบถ้วนเพื่อทำการบันทึกข้อมูล")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

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

	var visible string
	if postData.Visible == true {
		visible = "visible"
	} else {
		visible = "invisible"
	}

	// UPDATE PRODUCT
	sql := `UPDATE Product SET Name = ?, Description = ?, Price = ?, Cost = ?,  Barcode = ?,
			Status = ? , CategoryId = ?, Invisible =?, Quantity = ?, UpdateDate = ?
			WHERE ProductId = ? AND MerchantId = ?`
	_, err = db.Exec(sql, postData.Name, postData.Description, postData.Price,
		postData.Cost, postData.Barcode, postData.Status, postData.CategoryId,
		visible, postData.StockQuantity, config.GetCurrentDateTime(),
		pddata.ProductId, mcData.MerchantId,
	)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully update product",
		Data:       MerchantProductStruct.UpdateProductResponse{Success: true},
	}
	c.JSON(http.StatusOK, respond)
}
