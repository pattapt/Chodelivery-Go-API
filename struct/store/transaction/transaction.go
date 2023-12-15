package StoreTransactionStruct

import (
	StoreAddressStruct "gobasic/struct/store/address"
	StoreCartStruct "gobasic/struct/store/cart"
	StoreChatStruct "gobasic/struct/store/chat"
)

type TransactionListInfo struct {
	OrderId       int64                          `json:"orderId"`
	OrderToken    string                         `json:"orderToken"`
	PaymentMethod string                         `json:"paymentMethod"`
	TotalPay      float64                        `json:"totalPay"`
	Status        string                         `json:"status"`
	Note          string                         `json:"note"`
	DestinationId int64                          `json:"destinationId"`
	ChatId        int64                          `json:"chatId"`
	CreateDate    string                         `json:"createDate"`
	UpdateDate    string                         `json:"updateDate"`
	Merchant      MerchantDetails                `json:"merchant"`
	Chat          StoreChatStruct.ChatData       `json:"chat"`
	ChatProfile   StoreChatStruct.ChatProfile    `json:"chatProfile"`
	Destination   StoreAddressStruct.Destination `json:"destination"`
	Items         StoreCartStruct.CartData       `json:"items"`
}

type TransactionList struct {
	OrderId       int64           `json:"orderId"`
	OrderToken    string          `json:"orderToken"`
	PaymentMethod string          `json:"paymentMethod"`
	TotalPay      float64         `json:"totalPay"`
	Status        string          `json:"status"`
	Note          string          `json:"note"`
	ChatId        int64           `json:"chatId"`
	CreateDate    string          `json:"createDate"`
	UpdateDate    string          `json:"updateDate"`
	Merchant      MerchantDetails `json:"merchant"`
}

type MerchantDetails struct {
	MerchantId     int64  `json:"merchantId"`
	MerchantUUID   string `json:"merchantUUID"`
	OwnerSellerId  int64  `json:"ownerSellerId"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	PromptpayPhone string `json:"promptpayPhone"`
	Open           bool   `json:"open"`
	Visible        bool   `json:"visible"`
	ImageUrl       string `json:"imageUrl"`
	Address        string `json:"address"`
	Street         string `json:"street"`
	Building       string `json:"building"`
	District       string `json:"district"`
	Amphure        string `json:"amphure"`
	Province       string `json:"province"`
	ZipCode        string `json:"zipcode"`
}

type GetQRPromptPayResponse struct {
	Success    bool   `json:"success"`
	RawQRCode  string `json:"rawQrCode"`
	QRImageURL string `json:"qrImageUrl"`
}
