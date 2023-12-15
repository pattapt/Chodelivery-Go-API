package StoreV1Profile

import (
	"net/http"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	StoreAccountStruct "gobasic/struct/store/account"
)

func GetProfile(c *gin.Context) {

	db, err := database.ConnectDB()
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer db.Close()

	userID, exists := c.Get("userID")
	if !exists {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var account StoreAccountStruct.AccountData
	query := "SELECT * FROM Account WHERE AccountId = ?"
	err = db.QueryRow(query, userID).Scan(&account.AccountID, &account.AccountUUID, &account.Username, &account.Email, &account.Password, &account.ProfileImageURL, &account.RegisterDate, &account.LastLoginDate)
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
		Data: APIStructureMain.AccountProfile{
			AccountID:   account.AccountID,
			AccountUUID: account.AccountUUID,
			Username:    account.Username,
			Email:       account.Email,
			ProfileImageURL: func() string {
				if account.ProfileImageURL.Valid {
					return config.GetMainAPIURL() + account.ProfileImageURL.String
				}
				return ""
			}(),
			RegisterDate:  account.RegisterDate,
			LastLoginDate: account.LastLoginDate,
		},
	}
	c.JSON(http.StatusOK, respond)

}
