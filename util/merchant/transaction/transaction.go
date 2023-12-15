package UtilMerchant_Transaction

import (
	"fmt"
	"gobasic/config"
	"gobasic/database"
	MerchantTransactionStruct "gobasic/struct/merchant/transaction"
	StoreCartStruct "gobasic/struct/store/cart"
	UtilMerchant_Chat "gobasic/util/merchant/chat"
)

func GetTransactionDetail(OrderToken string, MerchantId int) (MerchantTransactionStruct.TransactionListDetail, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return MerchantTransactionStruct.TransactionListDetail{}, err
	}
	defer db.Close()
	query := `SELECT t.OrderId, t.OrderToken, t.PaymentMethod, t.TotalPay, t.Status,
		t.Note, t.ChatId, t.CreateDate, t.UpdateDate,
		cs.AccountId, cs.AccountUUID, cs.Username, cs.ProfileImageURL,
		ds.DestinationId, ds.DestinationToken, ds.Name AS DestinationName, ds.Phonenumber,
		ds.Address, ds.Street, ds.Building, ds.Note,
		d.name_th as district, a.name_th as amphure, p.name_th as province, d.zip_code as zipcode
		FROM Transaction t
		LEFT JOIN Merchant m ON t.MerchantId = m.MerchantId
		LEFT JOIN Account cs ON cs.AccountId = t.AccountId
		LEFT JOIN Destination ds ON t.DestinationId = ds.DestinationId
		LEFT JOIN districts d ON ds.distric = d.id
		LEFT JOIN amphures a ON d.amphure_id = a.id
		LEFT JOIN provinces p ON a.province_id = p.id
		WHERE t.OrderToken = ? AND m.MerchantId = ?`

	var tran MerchantTransactionStruct.TransactionListDetail
	err = db.QueryRow(query, OrderToken, MerchantId).Scan(
		&tran.OrderId,
		&tran.OrderToken,
		&tran.PaymentMethod,
		&tran.TotalPay,
		&tran.Status,
		&tran.Note,
		&tran.ChatId,
		&tran.CreateDate,
		&tran.UpdateDate,
		&tran.Customer.AccountId,
		&tran.Customer.AccountUUID,
		&tran.Customer.Username,
		&tran.Customer.ProfileImageURL,
		&tran.Destination.DestinationId,
		&tran.Destination.DestinationToken,
		&tran.Destination.DestinationName,
		&tran.Destination.PhoneNumber,
		&tran.Destination.Address,
		&tran.Destination.Street,
		&tran.Destination.Building,
		&tran.Destination.Note,
		&tran.Destination.District,
		&tran.Destination.Amphure,
		&tran.Destination.Province,
		&tran.Destination.ZipCode,
	)
	if err != nil {
		return MerchantTransactionStruct.TransactionListDetail{}, err
	}
	tran.Customer.ProfileImageURL = config.GetMainAPIURL() + tran.Customer.ProfileImageURL

	items, err := GetCartCheckOutByTransaction(int(tran.OrderId), int(tran.Customer.AccountId))
	if err != nil {
		return MerchantTransactionStruct.TransactionListDetail{}, err
	}
	tran.Items = items

	Chat, _ := UtilMerchant_Chat.GetChatInfoById(int(tran.ChatId))
	tran.Chat = Chat

	return tran, nil
}

func GetCartCheckOutByTransaction(TransactionId int, AccountId int) (StoreCartStruct.CartData, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return StoreCartStruct.CartData{}, err
	}
	defer db.Close()

	query := `SELECT c.CartId, c.CartToken, c.Amount, c.Status, c.CreateDate, c.UpdateDate,
			p.ProductId, p.ProductToken, p.Name, p.ImageUrl, p.Price, p.Price * c.Amount AS TotalPrice,
			CASE WHEN p.Status = 'available' THEN true ELSE false END AS Available,
			CASE WHEN p.Invisible = 'visible' THEN true ELSE false END AS Visible,
			c.Description AS Note
			FROM Cart c
			INNER JOIN Product p ON p.ProductId = c.ProductId
			WHERE c.AccountId = ? AND c.Status = 'checkout' AND c.OrderId = ?
			ORDER BY c.CartId DESC
    `
	rows, err := db.Query(query, AccountId, TransactionId)
	if err != nil {
		fmt.Print(err.Error())
		return StoreCartStruct.CartData{}, err
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
			&cart.Product.ImageURL,
			&cart.Product.Price,
			&cart.TotalPrice,
			&cart.Product.Available,
			&cart.Product.Visible,
			&cart.Note,
		)
		cart.Product.ImageURL = config.GetMainAPIURL() + cart.Product.ImageURL

		if err != nil {
			fmt.Print(err.Error())
			return StoreCartStruct.CartData{}, err
		}
		carts = append(carts, cart)
	}

	if len(carts) == 0 {
		fmt.Print(err.Error())
		return StoreCartStruct.CartData{}, err
	}

	return StoreCartStruct.CartData{
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
