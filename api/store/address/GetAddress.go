package StoreV1Address

import (
	"net/http"

	"github.com/gin-gonic/gin"

	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	StoreAddressStruct "gobasic/struct/store/address"
)

func GetAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
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

	query := `SELECT d.DestinationId, d.DestinationToken, d.Name, d.PhoneNumber, d.Status, d.CreateDate, d.UpdateDate,
				d.Address, d.Street, d.Building, dt.id as DistrictId, dt.name_th as district, a.name_th as amphure, p.name_th as province, 
				dt.zip_code as zipcode, d.Note
				FROM Destination d
				LEFT JOIN districts dt ON d.distric = dt.id
				LEFT JOIN amphures a ON dt.amphure_id = a.id
				LEFT JOIN provinces p ON a.province_id = p.id
				WHERE d.AccountId = ? AND d.Status = 'show' ORDER BY d.DestinationId DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer rows.Close()

	var addressList []StoreAddressStruct.Destination

	for rows.Next() {
		var addr StoreAddressStruct.Destination

		err := rows.Scan(
			&addr.DestinationId,
			&addr.DestinationToken,
			&addr.Name,
			&addr.PhoneNumber,
			&addr.Status,
			&addr.CreateDate,
			&addr.UpdateDate,
			&addr.Address.Address,
			&addr.Address.Street,
			&addr.Address.Building,
			&addr.Address.DistrictId,
			&addr.Address.District,
			&addr.Address.Amphure,
			&addr.Address.Province,
			&addr.Address.ZipCode,
			&addr.Note,
		)

		if err != nil {
			errorResponse := APIStructureMain.GeneralErrorMSG()
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		addressList = append(addressList, addr)
	}

	if len(addressList) == 0 {
		// Return an empty array as the response
		respond := APIStructureMain.Response{
			StatusCode: http.StatusOK,
			Status:     "S",
			Message:    "No address found",
			Data:       []StoreAddressStruct.Destination{},
		}
		c.JSON(http.StatusOK, respond)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully get address",
		Data:       addressList,
	}
	c.JSON(http.StatusOK, respond)
}
