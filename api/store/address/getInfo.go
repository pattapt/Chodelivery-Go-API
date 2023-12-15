package StoreV1Address

import (
	"net/http"

	"github.com/gin-gonic/gin"

	APIStructureMain "gobasic/struct"
	UtilStore_Address "gobasic/util/store/address"
)

func GetAddressInfo(c *gin.Context) {
	AddressToken := c.Param("AddressToken")
	userID, exists := c.Get("userID")
	if !exists {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	address, err := UtilStore_Address.GetAddressByToken(AddressToken, int(userID.(int)))
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully retrieved data",
		Data:       address,
	}
	c.JSON(http.StatusOK, respond)
}
