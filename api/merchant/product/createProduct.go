package MerchantV1Product

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
	MerchantProductStruct "gobasic/struct/merchant/product"
)

func CreateProduct(c *gin.Context) {
	var postData MerchantProductStruct.CreateroductPost

	if err := c.ShouldBindJSON(&postData); err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Name == "" || postData.Description == "" || postData.Price == 0 || postData.Cost == 0 || postData.Status == "" {
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

	ProductToken := config.GenerateRefreshToken()

	// UPDATE PRODUCT
	sql := `INSERT INTO Product (ProductToken, MerchantId, CategoryId, Name, Description, Price, Cost, Barcode, Status, Invisible, Quantity, CreateDate, UpdateDate, CreateIp) 
			VALUE (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	Result, err := db.Exec(sql, ProductToken, mcData.MerchantId, postData.CategoryId, postData.Name, postData.Description, postData.Price, postData.Cost, postData.Barcode, postData.Status, visible, postData.StockQuantity, config.GetCurrentDateTime(), config.GetCurrentDateTime(), c.ClientIP())
	if err != nil {
		fmt.Print(err)
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	ProductId, _ := Result.LastInsertId()

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully create product",
		Data:       MerchantProductStruct.CreateProductResponse{Success: true, ProductId: ProductId, ProductToken: ProductToken},
	}
	c.JSON(http.StatusOK, respond)
}
