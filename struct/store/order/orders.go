package StoreOrderStruct

type Destination struct {
	DestinationId    int     `json:"destinationId"`
	DestinationToken string  `json:"destinationToken"`
	Name             string  `json:"name"`
	PhoneNumber      string  `json:"phoneNumber"`
	Status           string  `json:"status"`
	CreateDate       string  `json:"createDate"`
	UpdateDate       string  `json:"updateDate"`
	Address          Address `json:"address"`
	Note             string  `json:"note"`
}

type Address struct {
	Address  string `json:"address"`
	Street   string `json:"street"`
	Building string `json:"building"`
	District string `json:"district"`
	Amphure  string `json:"amphure"`
	Province string `json:"province"`
	ZipCode  string `json:"zipcode"`
}

type CreateAddressPost struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
	Street      string `json:"street"`
	Building    string `json:"building"`
	District    int    `json:"district"`
	Amphure     int    `json:"amphure"`
	Province    int    `json:"province"`
	ZipCode     string `json:"zipcode"`
	Note        string `json:"note"`
}

type CreateAddressResponse struct {
	Success          bool   `json:"success"`
	DestinationId    int64  `json:"destinationId"`
	DestinationToken string `json:"destinationToken"`
}

type UpdateAddressPost struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Status      string `json:"status"`
	Address     string `json:"address"`
	Street      string `json:"street"`
	Building    string `json:"building"`
	District    int    `json:"district"`
	Amphure     int    `json:"amphure"`
	Province    int    `json:"province"`
	ZipCode     string `json:"zipcode"`
	Note        string `json:"note"`
}
