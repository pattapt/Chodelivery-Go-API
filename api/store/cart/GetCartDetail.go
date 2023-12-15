package StoreV1Cartwdad

import (
	"net/http"

	"github.com/gin-gonic/gin"

	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	StoreCartStruct "gobasic/struct/store/cart"
)

func GetCartDetail(c *gin.Context) {
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

	CartToken := c.Param("CartToken")
	if CartToken == "" {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}
	var cart StoreCartStruct.CartsResponse
	query := `SELECT c.CartId, c.CartToken, c.Amount, c.Status, c.CreateDate, c.UpdateDate,
				p.ProductId, p.ProductToken, p.Name, p.Description, p.ImageUrl, p.Price, p.Price * c.Amount AS TotalPrice,
				CASE WHEN p.Status = 'available' THEN true ELSE false END AS Available,
				CASE WHEN p.Invisible = 'visible' THEN true ELSE false END AS Visible,
				c.Description AS Note
				FROM Cart c
				INNER JOIN Product p ON p.ProductId = c.ProductId
				WHERE c.AccountId = ? AND c.Status = 'wait' AND c.CartToken = ?`
	err = db.QueryRow(query, userID, CartToken).Scan(
		&cart.CartID,
		&cart.CartToken,
		&cart.Amount,
		&cart.Status,
		&cart.CreateDate,
		&cart.UpdateDate,
		&cart.Product.ProductID,
		&cart.Product.ProductToken,
		&cart.Product.Name,
		&cart.Product.Description,
		&cart.Product.ImageURL,
		&cart.Product.Price,
		&cart.TotalPrice,
		&cart.Product.Available,
		&cart.Product.Visible,
		&cart.Note,
	)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully get cart",
		Data:       cart,
	}
	c.JSON(http.StatusOK, respond)
}
