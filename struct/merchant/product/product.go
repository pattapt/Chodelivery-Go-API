package MerchantProductStruct

type Product struct {
	ProductId     int     `json:"productId"`
	ProductToken  string  `json:"productToken"`
	Barcode       string  `json:"barcode"`
	CategoryId    int     `json:"categoryId"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	ImageUrl      string  `json:"imageUrl"`
	Price         float32 `json:"price"`
	Cost          float32 `json:"cost"`
	StockQuantity int     `json:"quantity"`
	Status        string  `json:"status"`
	Visible       bool    `json:"visible"`
	CreateDate    string  `json:"createDate"`
	UpdateDate    string  `json:"updateDate"`
}

type EditProductPost struct {
	ProductId     int     `json:"productId"`
	Barcode       string  `json:"barcode"`
	CategoryId    int     `json:"categoryId"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float32 `json:"price"`
	Cost          float32 `json:"cost"`
	StockQuantity int     `json:"quantity"`
	Status        string  `json:"status"`
	Visible       bool    `json:"visible"`
}

type UpdateProductResponse struct {
	Success bool `json:"success"`
}

type CreateroductPost struct {
	Barcode       string  `json:"barcode"`
	CategoryId    int     `json:"categoryId"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float32 `json:"price"`
	Cost          float32 `json:"cost"`
	StockQuantity int     `json:"quantity"`
	Status        string  `json:"status"`
	Visible       bool    `json:"visible"`
}

type CreateProductResponse struct {
	Success      bool   `json:"success"`
	ProductId    int64  `json:"productId"`
	ProductToken string `json:"productToken"`
}
