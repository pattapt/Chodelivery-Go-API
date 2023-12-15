package StoreV1Transaction

import (
	"net/http"

	"github.com/gin-gonic/gin"

	APIStructureMain "gobasic/struct"
	UtilStore_Transaction "gobasic/util/store/transaction"
)

func GetTransactionInfo(c *gin.Context) {
	TransactionToken := c.Param("TransactionToken")
	userID, exists := c.Get("userID")
	if !exists {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	transaction, err := UtilStore_Transaction.GetTransactionByTokenInfo(TransactionToken, int(userID.(int)))
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully retrieved data",
		Data:       transaction,
	}
	c.JSON(http.StatusOK, respond)
}
