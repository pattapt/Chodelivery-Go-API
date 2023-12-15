package StoreV1Address

import (
	"net/http"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	StoreAddressStruct "gobasic/struct/store/address"
)

func CreateAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var postData StoreAddressStruct.CreateAddressPost
	if err := c.ShouldBindJSON(&postData); err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if postData.Name == "" || postData.PhoneNumber == "" || postData.Address == "" || postData.Street == "" || postData.Building == "" || postData.District == 0 {
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

	DestinationToken := config.GenerateRefreshToken()
	// ADD ADDRESS
	sql := `INSERT INTO Destination (DestinationToken, AccountId, Name, PhoneNumber, Address, Street, Building, distric, Note, CreateDate, UpdateDate, CreateIp)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	Result, err := db.Exec(sql, DestinationToken, userID, postData.Name, postData.PhoneNumber, postData.Address, postData.Street, postData.Building, postData.District, postData.Note, config.GetCurrentDateTime(), config.GetCurrentDateTime(), c.ClientIP())
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	DestinationId, _ := Result.LastInsertId()

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully add new cart",
		Data: StoreAddressStruct.CreateAddressResponse{
			Success:          true,
			DestinationToken: DestinationToken,
			DestinationId:    DestinationId,
		},
	}
	c.JSON(http.StatusOK, respond)
}
