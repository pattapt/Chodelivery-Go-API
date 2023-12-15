package StoreCartStruct

type AddCartPost struct {
	ItemID int    `json:"itemId"`
	Amount int    `json:"amount"`
	Note   string `json:"note"`
}

type AddCartResponse struct {
	CartId    int64  `json:"CartId"`
	CartToken string `json:"cartToken"`
	Success   bool   `json:"success"`
}

type CartData struct {
	TotalPrice float64         `json:"totalPrice"`
	Amount     int             `json:"amount"`
	Cart       []CartsResponse `json:"cart"`
}

type CartDataV2 struct {
	TotalPrice float64           `json:"totalPrice"`
	Amount     int               `json:"amount"`
	MerchantID int               `json:"merchantId"`
	Cart       []CartsResponseV2 `json:"cart"`
}

type CartsResponse struct {
	CartID     int            `json:"cartId"`
	CartToken  string         `json:"cartToken"`
	Amount     int            `json:"amount"`
	TotalPrice float64        `json:"totalPrice"`
	Status     string         `json:"status"`
	CreateDate string         `json:"createDate"`
	UpdateDate string         `json:"updateDate"`
	Product    ProductDetails `json:"product"`
	Note       string         `json:"note"`
}

type ProductDetails struct {
	ProductID    int     `json:"productId"`
	ProductToken string  `json:"productToken"`
	Name         string  `json:"name"`
	ImageURL     string  `json:"imageUrl"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Available    bool    `json:"available"`
	Visible      bool    `json:"visible"`
}

type CartsResponseV2 struct {
	CartID     int              `json:"cartId"`
	CartToken  string           `json:"cartToken"`
	Amount     int              `json:"amount"`
	TotalPrice float64          `json:"totalPrice"`
	Status     string           `json:"status"`
	CreateDate string           `json:"createDate"`
	UpdateDate string           `json:"updateDate"`
	Product    ProductDetailsV2 `json:"product"`
	Note       string           `json:"note"`
}

type ProductDetailsV2 struct {
	ProductID    int     `json:"productId"`
	ProductToken string  `json:"productToken"`
	MerchantID   int     `json:"merchantId"`
	Name         string  `json:"name"`
	ImageURL     string  `json:"imageUrl"`
	Price        float64 `json:"price"`
	Available    bool    `json:"available"`
	Visible      bool    `json:"visible"`
}

type EditCartPost struct {
	CartId int    `json:"cartId"`
	Amount int    `json:"amount"`
	Note   string `json:"note"`
}

type EditCartResponse struct {
	CartId    int64  `json:"CartId"`
	CartToken string `json:"cartToken"`
	Success   bool   `json:"success"`
}

type CheckOutCartPost struct {
	PaymentMethod string `json:"paymentMethod"`
	DestinationId int    `json:"destinationId"`
	Note          string `json:"note"`
}

type CheckOutCartResponse struct {
	Success          bool     `json:"success"`
	TransactionId    int64    `json:"transactionId"`
	TransactionToken string   `json:"transactionToken"`
	QrCodePayment    string   `json:"qrCodePayment"`
	Chat             ChatData `json:"chat"`
}

type ChatData struct {
	ChatId    int64  `json:"chatId"`
	ChatToken string `json:"chatToken"`
}
