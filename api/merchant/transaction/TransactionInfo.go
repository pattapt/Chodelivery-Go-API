package MerchantV1Transaction

import (
	"net/http"

	"github.com/gin-gonic/gin"

	APIStructureMain "gobasic/struct"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
	UtilMerchant_Transaction "gobasic/util/merchant/transaction"
	UtilStore_Chat "gobasic/util/store/chat"
)

func GetTransactionDetail(c *gin.Context) {
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

	OrderToken := c.Param("OrderToken")
	if OrderToken == "" {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "ข้อมูลที่คุณระบุไม่ถูกต้อง โปรดตรวจสอบข้อมูลแล้วทำรายการใหม่อีกครั้ง")
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	trans, err := UtilMerchant_Transaction.GetTransactionDetail(OrderToken, mcData.MerchantId)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	sellerIdRaw, exists := c.Get("sellerId")
	if !exists {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	ChatProfile, _ := UtilStore_Chat.GetWhoIMI(int(trans.ChatId), sellerIdRaw, "seller")
	trans.ChatProfile = ChatProfile

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully get transaction",
		Data:       trans,
	}
	c.JSON(http.StatusOK, respond)
}
