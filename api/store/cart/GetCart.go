package StoreV1Cartwdad

import (
	"net/http"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	StoreCartStruct "gobasic/struct/store/cart"
)

func GetCart(c *gin.Context) {
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

	query := `SELECT c.CartId, c.CartToken, c.Amount, c.Status, c.CreateDate, c.UpdateDate,
				p.ProductId, p.ProductToken, p.Name, p.Description, p.ImageUrl, p.Price, p.Price * c.Amount AS TotalPrice,
				CASE WHEN p.Status = 'available' THEN true ELSE false END AS Available,
				CASE WHEN p.Invisible = 'visible' THEN true ELSE false END AS Visible,
				c.Description AS Note
				FROM Cart c
				INNER JOIN Product p ON p.ProductId = c.ProductId
				WHERE c.AccountId = ? AND c.Status = 'wait'
				ORDER BY c.CartId DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer rows.Close()

	var carts []StoreCartStruct.CartsResponse

	for rows.Next() {
		var cart StoreCartStruct.CartsResponse

		err := rows.Scan(
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
		cart.Product.ImageURL = config.GetMainAPIURL() + cart.Product.ImageURL
		carts = append(carts, cart)
	}

	if len(carts) == 0 {
		// Return an empty array as the response
		respond := APIStructureMain.Response{
			StatusCode: http.StatusOK,
			Status:     "S",
			Message:    "No cart found",
			Data: StoreCartStruct.CartData{
				TotalPrice: 0,
				Amount:     0,
				Cart:       []StoreCartStruct.CartsResponse{},
			},
		}
		c.JSON(http.StatusOK, respond)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully get cart",
		Data: StoreCartStruct.CartData{
			TotalPrice: func() float64 {
				var totalPrice float64
				for _, cart := range carts {
					totalPrice += cart.TotalPrice
				}
				return totalPrice
			}(),
			Amount: len(carts),
			Cart:   carts,
		},
	}
	c.JSON(http.StatusOK, respond)
}
