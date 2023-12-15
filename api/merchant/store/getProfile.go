package MerchantV1Store

import (
	"net/http"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	StoreAccountStruct "gobasic/struct/merchant/account"
)

func GetProfile(c *gin.Context) {

	db, err := database.ConnectDB()
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer db.Close()

	sellerId, exists := c.Get("sellerId")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var account StoreAccountStruct.SellerAccount
	query := "SELECT * FROM Seller WHERE SellerId = ?"
	err = db.QueryRow(query, sellerId).Scan(&account.SellerID, &account.SellerUUID, &account.Role, &account.Name, &account.LastName, &account.Email, &account.Password, &account.MerchantID, &account.ProfileImageURL, &account.RegisterDate, &account.LastLoginDate, &account.IP)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Get Profile
	respond := APIStructureMain.Response{
		StatusCode: 200,
		Status:     "S",
		Message:    "Successfully Get Data",
		Data: StoreAccountStruct.SellerProfile{
			SellerID:        account.SellerID,
			SellerUUID:      account.SellerUUID,
			Role:            account.Role,
			Name:            account.Name,
			LastName:        account.LastName,
			Email:           account.Email,
			MerchantID:      int(account.MerchantID.Int64),
			ProfileImageURL: config.GetMainAPIURL() + account.ProfileImageURL,
			RegisterDate:    account.RegisterDate,
			LastLoginDate:   account.LastLoginDate,
		},
	}
	c.JSON(http.StatusOK, respond)

}
