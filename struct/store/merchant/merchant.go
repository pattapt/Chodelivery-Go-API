package StoreMerchantStruct

type Category struct {
	CategoryId    int    `json:"categoryId"`
	CategoryToken string `json:"categoryToken"`
	MerchantId    int    `json:"merchantId"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ImageUrl      string `json:"imageUrl"`
	Status        string `json:"status"`
	Invisible     string `json:"invisible"`
	CreateDate    string `json:"createDate"`
	UpdateDate    string `json:"updateDate"`
	CreateIp      string `json:"-"`
}

type CategoryWithProducts struct {
	CategoryId    int       `json:"categoryId"`
	CategoryToken string    `json:"categoryToken"`
	MerchantId    int       `json:"merchantId"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ImageUrl      string    `json:"imageUrl"`
	Status        string    `json:"status"`
	Invisible     string    `json:"invisible"`
	CreateDate    string    `json:"createDate"`
	UpdateDate    string    `json:"updateDate"`
	Products      []Product `json:"products"`
}

type Product struct {
	ProductId         int     `json:"productId"`
	ProductToken      string  `json:"productToken"`
	Barcode           string  `json:"barcode"`
	ProductCategoryId int     `json:"productCategoryId"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	ImageUrl          string  `json:"imageUrl"`
	Price             float32 `json:"price"`
	StockQuantity     int     `json:"quantity"`
	Status            string  `json:"status"`
	CreateDate        string  `json:"createDate"`
	UpdateDate        string  `json:"updateDate"`
}

type MerchantList struct {
	PageId    int              `json:"pageId"`
	TotalPage int              `json:"totalPage"`
	TotalItem int              `json:"totalItem"`
	Merchants []MerchantDetail `json:"merchants"`
}

type MerchantDetail struct {
	MerchantId   int    `json:"id"`
	MerchantUUID string `json:"uuid"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Open         bool   `json:"open"`
	Visible      bool   `json:"visible"`
	ImageUrl     string `json:"image_url"`
	Address      struct {
		Address  string `json:"address"`
		Street   string `json:"street"`
		Building string `json:"building"`
		District string `json:"district"`
		Amphure  string `json:"amphure"`
		Province string `json:"province"`
		ZipCode  string `json:"zipcode"`
	}
	CreateDate string `json:"CreateDate"`
	UpdateDate string `json:"UpdateDate"`
}

type MerchantDetailV2 struct {
	MerchantId   int    `json:"id"`
	MerchantUUID string `json:"uuid"`
	OwnerID      string `json:"owner_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Open         bool   `json:"open"`
	Visible      bool   `json:"visible"`
	ImageUrl     string `json:"image_url"`
	Address      struct {
		Address  string `json:"address"`
		Street   string `json:"street"`
		Building string `json:"building"`
		District string `json:"district"`
		Amphure  string `json:"amphure"`
		Province string `json:"province"`
		ZipCode  string `json:"zipcode"`
	}
	CreateDate string `json:"CreateDate"`
	UpdateDate string `json:"UpdateDate"`
}
