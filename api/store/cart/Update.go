package StoreV1Cartwdad

import (
	"net/http"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	StoreCartStruct "gobasic/struct/store/cart"

	UtilStore_Cart "gobasic/util/store/cart"
)

func UpdateCart(c *gin.Context) {
	var postData StoreCartStruct.EditCartPost
	if err := c.ShouldBindJSON(&postData); err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.CartId == 0 {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	cart, err := UtilStore_Cart.GetCartById(postData.CartId)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if cart.Status != "wait" {
		errorResponse := APIStructureMain.ErrorRespondAuth("ไม่สามารถแก้ไขข้อมูลได้", "ระบบไม่สามารถแก้ไขได้เนื่องจากไม่มีสินค้านี้ในตระกร้าของคุณแล้ว")
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

	// Update Cart
	status := "wait"
	if postData.Amount == 0 {
		status = "remove"
	}
	sql := `UPDATE Cart SET Amount = ?, Description = ?, Status = ?, UpdateDate = ? WHERE CartId = ? AND AccountId = ? AND Status = 'wait'`
	_, err = db.Exec(sql, postData.Amount, postData.Note, status, config.GetCurrentDateTime(), cart.CartID, userID)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully update cart",
		Data: StoreCartStruct.EditCartResponse{
			CartId:  int64(cart.CartID),
			Success: true,
		},
	}
	c.JSON(http.StatusOK, respond)
}
