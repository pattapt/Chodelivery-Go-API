package MerchantV1Product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
	MerchantProductStruct "gobasic/struct/merchant/product"
)

func GetProduct(c *gin.Context) {
	merchantData, exists := c.Get("merchantData")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	mcData, ok := merchantData.(MerchantMerchantStruct.MerchantDetail)
	if !ok {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

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

	query := `SELECT p.ProductId, p.ProductToken, p.Barcode, p.CategoryId, 
				p.Name, p.Description, p.ImageUrl, p.Price, p.Cost, p.Quantity, 
				p.Status,
				CASE WHEN p.Invisible = 'visible' THEN true ELSE false END AS visible,
				p.CreateDate, p.UpdateDate
				FROM Product p WHERE p.MerchantId = ? 
				ORDER BY p.ProductId DESC LIMIT ?, ?
				`
	rows, err := db.Query(query, mcData.MerchantId, Start, limit)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer rows.Close()

	var products []MerchantProductStruct.Product

	for rows.Next() {
		var pd MerchantProductStruct.Product

		err := rows.Scan(
			&pd.ProductId,
			&pd.ProductToken,
			&pd.Barcode,
			&pd.CategoryId,
			&pd.Name,
			&pd.Description,
			&pd.ImageUrl,
			&pd.Price,
			&pd.Cost,
			&pd.StockQuantity,
			&pd.Status,
			&pd.Visible,
			&pd.CreateDate,
			&pd.UpdateDate,
		)
		pd.ImageUrl = config.GetMainAPIURL() + pd.ImageUrl

		if err != nil {
			// fmt.Print(err)
			errorResponse := APIStructureMain.GeneralErrorMSG()
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		products = append(products, pd)
	}

	if len(products) == 0 {
		// Return an empty array as the response
		respond := APIStructureMain.Response{
			StatusCode: http.StatusOK,
			Status:     "S",
			Message:    "No product found",
			Data:       []MerchantProductStruct.Product{},
		}
		c.JSON(http.StatusOK, respond)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully get product",
		Data:       products,
	}
	c.JSON(http.StatusOK, respond)
}
