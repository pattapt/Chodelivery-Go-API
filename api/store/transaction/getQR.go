package StoreV1Transaction

import (
	"bytes"
	"image"
	"image/png"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	promptpayqr "github.com/kazekim/promptpay-qr-go"

	APIStructureMain "gobasic/struct"
	UtilStore_Transaction "gobasic/util/store/transaction"
	// pp "github.com/Frontware/promptpay"
)

func GetQRCode(c *gin.Context) {
	TransactionToken := c.Param("TransactionToken")
	userID, exists := c.Get("userID")
	if !exists {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	transaction, err := UtilStore_Transaction.GetTransactionByToken(TransactionToken, int(userID.(int)))
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// payment := pp.PromptPay{
	// 	PromptPayID: "0105540087061",
	// 	Amount:      transaction.TotalPay,
	// }
	// qrcode, _ := payment.Gen()

	target := transaction.Merchant.PromptpayPhone
	amount := transaction.TotalPay

	qr, err := promptpayqr.QRForTargetWithAmount(target, strconv.FormatFloat(amount, 'f', 2, 64))
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Convert the QR code to an image
	img, _, err := image.Decode(bytes.NewReader(*qr))
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Return the QR code image as a response
	c.Header("Content-Type", "image/png")
	c.Status(http.StatusOK)
	err = png.Encode(c.Writer, img)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
}
