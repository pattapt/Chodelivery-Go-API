package StoreV1Cartwdad

import (
	"net/http"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	StoreCartStruct "gobasic/struct/store/cart"

	UtilStore_Product "gobasic/util/store/product"
)

func AddCart(c *gin.Context) {
	var postData StoreCartStruct.AddCartPost
	if err := c.ShouldBindJSON(&postData); err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.ItemID == 0 {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Amount < 0 {
		errorResponse := APIStructureMain.ErrorRespondAuth("ไม่สามารถทำรายการได้", "จำนวนสินค้าไม่ถูกต้อง")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	product, err := UtilStore_Product.GetProductById(postData.ItemID)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
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

	CartToken := config.GenerateRefreshToken()
	// ADD CART
	sql := `INSERT INTO Cart (CartToken, ProductId, AccountId, Amount, Description, CreateDate, UpdateDate, CreateIp) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	Result, err := db.Exec(sql, CartToken, product.ProductId, userID, postData.Amount, postData.Note, config.GetCurrentDateTime(), config.GetCurrentDateTime(), c.ClientIP())
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	CartId, err := Result.LastInsertId()

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully add new cart",
		Data: StoreCartStruct.AddCartResponse{
			Success:   true,
			CartToken: CartToken,
			CartId:    CartId,
		},
	}
	c.JSON(http.StatusOK, respond)
}
