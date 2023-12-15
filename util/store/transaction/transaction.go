package UtilStore_Transaction

import (
	"fmt"
	"gobasic/config"
	"gobasic/database"
	StoreCartStruct "gobasic/struct/store/cart"
	StoreTransactionStruct "gobasic/struct/store/transaction"
	UtilStore_Address "gobasic/util/store/address"
	UtilStore_Chat "gobasic/util/store/chat"
)

func GetTransactionByToken(TransactionToken string, AccountId int) (StoreTransactionStruct.TransactionList, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return StoreTransactionStruct.TransactionList{}, err
	}
	defer db.Close()
	var tran StoreTransactionStruct.TransactionList
	query := `
		SELECT t.OrderId, t.OrderToken, t.PaymentMethod, t.TotalPay, t.Status, 
		t.Note, t.ChatId, t.CreateDate, t.UpdateDate,
		m.MerchantId, m.MerchantUUID, m.OwnerSellerId, m.Name, m.Description, m.PromptpayPhone,
		CASE WHEN m.Status = 'open' THEN true ELSE false END AS Open,
		CASE WHEN m.Visible = 'visible' THEN true ELSE false END AS Visible,
		m.ImageUrl, m.Address, m.Street, m.Building, 
		d.name_th as district, a.name_th as amphure, p.name_th as province, d.zip_code as zipcode
		FROM Transaction t
		LEFT JOIN Merchant m ON t.MerchantId = m.MerchantId
		LEFT JOIN districts d ON m.distric = d.id
		LEFT JOIN amphures a ON d.amphure_id = a.id
		LEFT JOIN provinces p ON a.province_id = p.id
		WHERE t.OrderToken = ? AND t.AccountId = ?
    `
	err = db.QueryRow(query, TransactionToken, AccountId).Scan(
		&tran.OrderId,
		&tran.OrderToken,
		&tran.PaymentMethod,
		&tran.TotalPay,
		&tran.Status,
		&tran.Note,
		&tran.ChatId,
		&tran.CreateDate,
		&tran.UpdateDate,
		&tran.Merchant.MerchantId,
		&tran.Merchant.MerchantUUID,
		&tran.Merchant.OwnerSellerId,
		&tran.Merchant.Name,
		&tran.Merchant.Description,
		&tran.Merchant.PromptpayPhone,
		&tran.Merchant.Open,
		&tran.Merchant.Visible,
		&tran.Merchant.ImageUrl,
		&tran.Merchant.Address,
		&tran.Merchant.Street,
		&tran.Merchant.Building,
		&tran.Merchant.District,
		&tran.Merchant.Amphure,
		&tran.Merchant.Province,
		&tran.Merchant.ZipCode,
	)
	tran.Merchant.ImageUrl = config.GetMainAPIURL() + tran.Merchant.ImageUrl
	if err != nil {
		fmt.Print(err.Error())
		return StoreTransactionStruct.TransactionList{}, err
	}

	return tran, nil
}

func GetTransactionByTokenInfo(TransactionToken string, AccountId int) (StoreTransactionStruct.TransactionListInfo, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return StoreTransactionStruct.TransactionListInfo{}, err
	}
	defer db.Close()
	var tran StoreTransactionStruct.TransactionListInfo
	query := `
		SELECT t.OrderId, t.OrderToken, t.PaymentMethod, t.TotalPay, t.Status, 
		t.Note, t.DestinationId, t.ChatId, t.CreateDate, t.UpdateDate,
		m.MerchantId, m.MerchantUUID, m.OwnerSellerId, m.Name, m.Description,
		CASE WHEN m.Status = 'open' THEN true ELSE false END AS Open,
		CASE WHEN m.Visible = 'visible' THEN true ELSE false END AS Visible,
		m.ImageUrl, m.Address, m.Street, m.Building, 
		d.name_th as district, a.name_th as amphure, p.name_th as province, d.zip_code as zipcode
		FROM Transaction t
		LEFT JOIN Merchant m ON t.MerchantId = m.MerchantId
		LEFT JOIN districts d ON m.distric = d.id
		LEFT JOIN amphures a ON d.amphure_id = a.id
		LEFT JOIN provinces p ON a.province_id = p.id
		WHERE t.OrderToken = ? AND t.AccountId = ?
    `
	err = db.QueryRow(query, TransactionToken, AccountId).Scan(
		&tran.OrderId,
		&tran.OrderToken,
		&tran.PaymentMethod,
		&tran.TotalPay,
		&tran.Status,
		&tran.Note,
		&tran.DestinationId,
		&tran.ChatId,
		&tran.CreateDate,
		&tran.UpdateDate,
		&tran.Merchant.MerchantId,
		&tran.Merchant.MerchantUUID,
		&tran.Merchant.OwnerSellerId,
		&tran.Merchant.Name,
		&tran.Merchant.Description,
		&tran.Merchant.Open,
		&tran.Merchant.Visible,
		&tran.Merchant.ImageUrl,
		&tran.Merchant.Address,
		&tran.Merchant.Street,
		&tran.Merchant.Building,
		&tran.Merchant.District,
		&tran.Merchant.Amphure,
		&tran.Merchant.Province,
		&tran.Merchant.ZipCode,
	)
	tran.Merchant.ImageUrl = config.GetMainAPIURL() + tran.Merchant.ImageUrl
	if err != nil {
		fmt.Print(err.Error())
		return StoreTransactionStruct.TransactionListInfo{}, err
	}
	Chat, _ := UtilStore_Chat.GetChatInfoById(int(tran.ChatId))
	tran.Chat = Chat

	CartCheckOut, _ := GetCartCheckOutByTransaction(int(tran.OrderId), AccountId)
	tran.Items = CartCheckOut

	Destination, _ := UtilStore_Address.GetAddressById(int(tran.DestinationId), AccountId)
	tran.Destination = Destination

	ChatProfile, _ := UtilStore_Chat.GetWhoIMI(int(tran.ChatId), AccountId, "customer")
	tran.ChatProfile = ChatProfile

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
