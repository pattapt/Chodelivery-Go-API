package MerchantUtilStruct

type THProvince struct {
	District District `json:"district"`
	Amphure  Amphure  `json:"amphure"`
	Province Province `json:"province"`
	ZipCode  string   `json:"zipcode"`
}

type District struct {
	DistrictId int    `json:"district_id"`
	NameTh     string `json:"name_th"`
	NameEn     string `json:"name_en"`
}

type Province struct {
	ProvinceId int    `json:"province_id"`
	NameTh     string `json:"name_th"`
	NameEn     string `json:"name_en"`
}

type Amphure struct {
	AmphureId int    `json:"amphure_id"`
	NameTh    string `json:"name_th"`
	NameEn    string `json:"name_en"`
}
