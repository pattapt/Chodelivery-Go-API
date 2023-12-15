package StoreV1Address

import (
	"net/http"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	StoreAddressStruct "gobasic/struct/store/address"
	UtilStore_Address "gobasic/util/store/address"
)

func UpdateAddress(c *gin.Context) {
	AddressToken := c.Param("AddressToken")
	userID, exists := c.Get("userID")
	if !exists {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var postData StoreAddressStruct.UpdateAddressPost
	if err := c.ShouldBindJSON(&postData); err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Name == "" || postData.PhoneNumber == "" || postData.Address == "" || postData.Street == "" || postData.Building == "" || postData.Status == "" || postData.District == 0 {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	address, _ := UtilStore_Address.GetAddressByToken(AddressToken, int(userID.(int)))
	if address.Status != "show" {
		errorResponse := APIStructureMain.ErrorRespondAuth("ไม่สามารถแก้ไขข้อมูลได้", "ท่านไม่สามารถทำการแก้ไขข้อมูลได้แล้ว เนื่องจากไม่มีข้อมูลที่อยู่นี้แล้ว")
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

	// UPDATE ADDRESS
	AddressId := address.DestinationId
	sql := ""
	if postData.Status == "delete" {
		sql = `UPDATE Destination SET Status = 'delete', UpdateDate = ? WHERE DestinationId = ? AND AccountId = ?`
		_, err = db.Exec(sql, config.GetCurrentDateTime(), AddressId, userID)
	} else {
		sql = `UPDATE Destination SET Name = ?, PhoneNumber = ?, Address = ?, Street = ?, 
			Building = ?, distric = ?, Note = ?, UpdateDate = ? 
			WHERE DestinationId = ? AND AccountId = ?`
		_, err = db.Exec(sql, postData.Name, postData.PhoneNumber, postData.Address, postData.Street, postData.Building, postData.District, postData.Note, config.GetCurrentDateTime(), AddressId, userID)
	}

	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully update address",
		Data: StoreAddressStruct.CreateAddressResponse{
			Success:          true,
			DestinationId:    int64(address.DestinationId),
			DestinationToken: address.DestinationToken,
		},
	}
	c.JSON(http.StatusOK, respond)
}
