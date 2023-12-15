package MerchantProfileStruct

// import "database/sql"

type MerchantProfile struct {
	MerchantID    int    `json:"merchant_id"`
	MerchantUUID  string `json:"merchant_uuid"`
	OwnerSellerID int    `json:"owner_seller_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	ImageURL      string `json:"image_url"`
	Visible       string `json:"visible"`
	Address       string `json:"address"`
	Street        string `json:"street"`
	Building      string `json:"building"`
	District      int    `json:"district"`
	CreateDate    string `json:"create_date"`
	UpdateDate    string `json:"update_date"`
	CreateIP      string `json:"create_ip"`
}

type CreateMerchantRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PromptpayId string `json:"promptpayId"`
	Address     string `json:"address"`
	Street      string `json:"street"`
	Building    string `json:"building"`
	District    int    `json:"district"`
	Amphures    int    `json:"amphures"`
	Provinces   int    `json:"provinces"`
}

type MerchantList struct {
	PageId    int              `json:"pageId"`
	TotalPage int              `json:"totalPage"`
	TotalItem int              `json:"totalItem"`
	Merchants []MerchantDetail `json:"merchants"`
}

type MerchantDetail struct {
	MerchantId     int    `json:"id"`
	MerchantUUID   string `json:"uuid"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	PromptPayPhone string `json:"promptPayPhone"`
	Open           bool   `json:"open"`
	Visible        bool   `json:"visible"`
	ImageUrl       string `json:"image_url"`
	Address        struct {
		Address    string `json:"address"`
		Street     string `json:"street"`
		Building   string `json:"building"`
		District   string `json:"district"`
		DistrictId int    `json:"district_id"`
		Amphure    string `json:"amphure"`
		Province   string `json:"province"`
		ZipCode    string `json:"zipcode"`
	}
	CreateDate string `json:"CreateDate"`
	UpdateDate string `json:"UpdateDate"`
}

type EditMerchantProfilePost struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	PromptpayPhone string `json:"promptPayPhone"`
	Open           bool   `json:"open"`
	Visible        bool   `json:"visible"`
	Address        string `json:"address"`
	Street         string `json:"street"`
	Building       string `json:"building"`
	District       int    `json:"district"`
}

type UpdateProfileResponse struct {
	Success bool `json:"success"`
}

type MerchantSummary struct {
	TotalOrder    int     `json:"totalOrder"`
	TotalAmount   float64 `json:"totalAmount"`
	LastWeekOrder int     `json:"lastWeekOrder"`
	LastWeekSales float64 `json:"lastWeekSales"`
}
