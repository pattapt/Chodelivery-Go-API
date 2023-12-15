package StoreV1Merchant

import (
	"net/http"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"

	// StoreAccountStruct "gobasic/struct/store/account"
	StoreMerchantStruct "gobasic/struct/store/merchant"
)

func GetAllCategories(c *gin.Context) {
	MerchantId, exists := c.Get("merchantId")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
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

	query := `SELECT c.CategoryId, c.CategoryToken, c.MerchantId, c.Name, c.Description, c.ImageUrl, c.Status, c.Invisible, c.CreateDate, c.UpdateDate 
				FROM Category c 
				WHERE (c.MerchantId = ? OR c.MerchantId = 0)
				AND 0 != (SELECT COUNT(*) FROM Product p WHERE p.CategoryId = c.CategoryId AND p.Invisible = 'visible')`
	rows, err := db.Query(query, MerchantId)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer rows.Close()

	var categories []StoreMerchantStruct.Category

	for rows.Next() {
		var category StoreMerchantStruct.Category
		err := rows.Scan(
			&category.CategoryId,
			&category.CategoryToken,
			&category.MerchantId,
			&category.Name,
			&category.Description,
			&category.ImageUrl,
			&category.Status,
			&category.Invisible,
			&category.CreateDate,
			&category.UpdateDate,
		)
		if err != nil {
			errorResponse := APIStructureMain.GeneralErrorMSG()
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		category.ImageUrl = config.GetMainAPIURL() + category.ImageUrl

		categories = append(categories, category)
	}

	if len(categories) == 0 {
		// Return an empty array as the response
		respond := APIStructureMain.Response{
			StatusCode: http.StatusOK,
			Status:     "S",
			Message:    "No category found for the merchant",
			Data:       []StoreMerchantStruct.Category{},
		}
		c.JSON(http.StatusOK, respond)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully retrieved categories",
		Data:       categories,
	}
	c.JSON(http.StatusOK, respond)
}
