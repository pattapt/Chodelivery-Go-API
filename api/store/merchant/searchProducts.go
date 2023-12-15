package StoreV1Merchant

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"

	// StoreAccountStruct "gobasic/struct/store/account"
	StoreMerchantStruct "gobasic/struct/store/merchant"
)

func SearchProducts(c *gin.Context) {
	// Extract category ID from the URL parameters
	MerchantId, exists := c.Get("merchantId")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	keywords := c.DefaultQuery("keywords", "")
	limitStr := c.DefaultQuery("limit", "50")
	pageStr := c.DefaultQuery("page", "1")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 50
	}

	Start := (page - 1) * limit
	db, err := database.ConnectDB()
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer db.Close()

	// Query to get products for a specific category
	query := `
		SELECT 
			ProductId,
			ProductToken,
			Barcode,
			CategoryId AS ProductCategoryId,
			Name AS ProductName,
			Description AS ProductDescription,
			ImageUrl AS ProductImageUrl,
			Price AS ProductPrice,
			Status AS ProductStatus,
			CreateDate AS ProductCreateDate,
			UpdateDate AS ProductUpdateDate
		FROM Product
		WHERE Name LIKE ? AND MerchantId = ? AND Invisible = 'visible'
		ORDER BY CreateDate DESC
		LIMIT ?, ?
	`
	searchString := "%" + keywords + "%"
	rows, err := db.Query(query, searchString, MerchantId, Start, limit)
	if err != nil {
		fmt.Println(err)
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer rows.Close()

	var products []StoreMerchantStruct.Product

	for rows.Next() {
		var product StoreMerchantStruct.Product

		err := rows.Scan(
			&product.ProductId,
			&product.ProductToken,
			&product.Barcode,
			&product.ProductCategoryId,
			&product.Name,
			&product.Description,
			&product.ImageUrl,
			&product.Price,
			&product.Status,
			&product.CreateDate,
			&product.UpdateDate,
		)

		if err != nil {
			errorResponse := APIStructureMain.GeneralErrorMSG()
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		product.ImageUrl = config.GetMainAPIURL() + product.ImageUrl
		products = append(products, product)
	}

	if len(products) == 0 {
		// Return an empty array as the response
		respond := APIStructureMain.Response{
			StatusCode: http.StatusOK,
			Status:     "S",
			Message:    "No products found for the category",
			Data:       []StoreMerchantStruct.Product{},
		}
		c.JSON(http.StatusOK, respond)
		return
	}

	// Respond with the retrieved products
	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully retrieved products for the category",
		Data:       products,
	}
	c.JSON(http.StatusOK, respond)
}
