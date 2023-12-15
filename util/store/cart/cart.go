package UtilStore_Cart

import (
	"fmt"
	"gobasic/database"
	StoreCartStruct "gobasic/struct/store/cart"
)

func GetCartById(CartId int) (StoreCartStruct.CartsResponse, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return StoreCartStruct.CartsResponse{}, err
	}
	defer db.Close()

	var cart StoreCartStruct.CartsResponse
	query := `
		SELECT c.CartId, c.CartToken, c.Amount, c.Status, c.CreateDate, c.UpdateDate,
		p.ProductId, p.ProductToken, p.Name, p.ImageUrl, p.Price, p.Price * c.Amount AS TotalPrice,
		CASE WHEN p.Status = 'available' THEN true ELSE false END AS Available,
		CASE WHEN p.Invisible = 'visible' THEN true ELSE false END AS Visible,
		c.Description AS Note
		FROM Cart c
		INNER JOIN Product p ON p.ProductId = c.ProductId
		WHERE c.CartId = ?
    `
	err = db.QueryRow(query, CartId).Scan(
		&cart.CartID,
		&cart.CartToken,
		&cart.Amount,
		&cart.Status,
		&cart.CreateDate,
		&cart.UpdateDate,
		&cart.Product.ProductID,
		&cart.Product.ProductToken,
		&cart.Product.Name,
		&cart.Product.ImageURL,
		&cart.Product.Price,
		&cart.TotalPrice,
		&cart.Product.Available,
		&cart.Product.Visible,
		&cart.Note,
	)
	if err != nil {
		fmt.Print(err.Error())
		return StoreCartStruct.CartsResponse{}, err
	}

	return cart, nil
}

func GetAllCart(AccountId int) (StoreCartStruct.CartDataV2, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return StoreCartStruct.CartDataV2{}, err
	}
	defer db.Close()

	query := `SELECT c.CartId, c.CartToken, c.Amount, c.Status, c.CreateDate, c.UpdateDate,
			p.ProductId, p.ProductToken, p.MerchantId, p.Name, p.ImageUrl, p.Price, p.Price * c.Amount AS TotalPrice,
			CASE WHEN p.Status = 'available' THEN true ELSE false END AS Available,
			CASE WHEN p.Invisible = 'visible' THEN true ELSE false END AS Visible,
			c.Description AS Note
			FROM Cart c
			INNER JOIN Product p ON p.ProductId = c.ProductId
			WHERE c.AccountId = ? AND c.Status = 'wait'
			ORDER BY c.CartId DESC
    `
	rows, err := db.Query(query, AccountId)
	if err != nil {
		fmt.Print(err.Error())
		return StoreCartStruct.CartDataV2{}, err
	}
	defer rows.Close()

	var carts []StoreCartStruct.CartsResponseV2

	for rows.Next() {
		var cart StoreCartStruct.CartsResponseV2

		err := rows.Scan(
			&cart.CartID,
			&cart.CartToken,
			&cart.Amount,
			&cart.Status,
			&cart.CreateDate,
			&cart.UpdateDate,
			&cart.Product.ProductID,
			&cart.Product.ProductToken,
			&cart.Product.MerchantID,
			&cart.Product.Name,
			&cart.Product.ImageURL,
			&cart.Product.Price,
			&cart.TotalPrice,
			&cart.Product.Available,
			&cart.Product.Visible,
			&cart.Note,
		)

		if err != nil {
			fmt.Print(err.Error())
			return StoreCartStruct.CartDataV2{}, err
		}
		carts = append(carts, cart)
	}

	if len(carts) == 0 {
		// fmt.Print(err.Error())
		return StoreCartStruct.CartDataV2{}, err
	}

	return StoreCartStruct.CartDataV2{
		MerchantID: carts[0].Product.MerchantID,
		TotalPrice: func() float64 {
			var totalPrice float64
			for _, cart := range carts {
				totalPrice += cart.TotalPrice
			}
			return totalPrice
		}(),
		Amount: len(carts),
		Cart:   carts,
	}, nil
}
