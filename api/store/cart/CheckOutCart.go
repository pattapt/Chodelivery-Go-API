package StoreV1Cartwdad

import (
	"net/http"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	StoreCartStruct "gobasic/struct/store/cart"
	UtilStore_Cart "gobasic/util/store/cart"
	UtilStore_Product "gobasic/util/store/product"
)

func CheckOutCart(c *gin.Context) {
	var postData StoreCartStruct.CheckOutCartPost
	if err := c.ShouldBindJSON(&postData); err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.PaymentMethod == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("ไม่สามารถทำรายการได้", "กรุณาระบุวิธีการชำระเงิน")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.DestinationId == 0 {
		errorResponse := APIStructureMain.ErrorRespondAuth("ไม่สามารถทำรายการได้", "กรุณาเลือกที่อยู่สำหรับจัดส่งสินค้า")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	Carts, err := UtilStore_Cart.GetAllCart(userID.(int))
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if len(Carts.Cart) == 0 {
		errorResponse := APIStructureMain.ErrorRespondAuth("ไม่สามารถทำรายการได้", "ไม่พบสินค้าในตะกร้า เพิ่มสินค้าลงในตระกร้าก่อนแล้วทำรายการอีกครั้ง")
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

	MerchantData, _ := UtilStore_Product.GetMerchantInfoById(Carts.MerchantID)

	ChatToken := config.GenerateRefreshToken()
	sql := `INSERT INTO Chat (ChatToken, Status, CreateDate, LastTalkDate)
			VALUES (?, ?, ?, ?)`
	Result, _ := db.Exec(sql, ChatToken, "open", config.GetCurrentDateTime(), config.GetCurrentDateTime())
	ChatId, err := Result.LastInsertId()

	TransactionToken := config.GenerateRefreshToken()
	// ADD CART
	sql = `INSERT INTO Transaction (OrderToken, MerchantId, AccountId, PaymentMethod, TotalPay, Note, DestinationId, ChatId, CreateDate, UpdateDate, CreateIp)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	Result, err = db.Exec(sql, TransactionToken, Carts.MerchantID, userID, postData.PaymentMethod, Carts.TotalPrice, postData.Note, postData.DestinationId, ChatId, config.GetCurrentDateTime(), config.GetCurrentDateTime(), c.ClientIP())
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	TransactionId, err := Result.LastInsertId()

	// UPDATE ALL CART TO CHECK OUT
	sql = `UPDATE Cart SET Status = 'checkout', OrderId = ? WHERE AccountId = ? AND Status = 'wait'`
	_, err = db.Exec(sql, TransactionId, userID)

	// ADD MERCHANT OWNER AND CUSTOMER TO CHAT
	sql = `INSERT INTO ChatMember (MemberUUID, ChatId, AccountType, AccountId, Status)
			VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(sql, config.GenerateRefreshToken(), ChatId, "customer", userID, "join")

	// ต้องดึงข้อมูลจากตาราง Merchant มาใส่โดยจะใส่เจ้าของไว้ก่อน
	OwnerId := MerchantData.OwnerID
	sql = `INSERT INTO ChatMember (MemberUUID, ChatId, AccountType, AccountId, Status)
			VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(sql, config.GenerateRefreshToken(), ChatId, "seller", OwnerId, "join")

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully checkout cart",
		Data: StoreCartStruct.CheckOutCartResponse{
			Success:          true,
			TransactionId:    TransactionId,
			TransactionToken: TransactionToken,
			QrCodePayment:    config.GetMainAPIURL() + "/api/store/v1/transaction/" + TransactionToken + "/QR",
			Chat: StoreCartStruct.ChatData{
				ChatId:    int64(ChatId),
				ChatToken: config.GenerateRefreshToken(),
			},
		},
	}
	c.JSON(http.StatusOK, respond)
}
