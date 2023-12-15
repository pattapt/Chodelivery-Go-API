package MerchantAccountStruct

import "database/sql"

type RegisterAccountPost struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}

type LoginAccountPost struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type SellerAccount struct {
	SellerID        int           `json:"seller_id"`
	SellerUUID      string        `json:"seller_uuid"`
	Role            string        `json:"role"`
	Name            string        `json:"name"`
	LastName        string        `json:"last_name"`
	Email           string        `json:"email"`
	Password        string        `json:"password"`
	MerchantID      sql.NullInt64 `json:"merchant_id"`
	ProfileImageURL string        `json:"profile_image_url"`
	RegisterDate    string        `json:"register_date"`
	LastLoginDate   string        `json:"last_login_date"`
	IP              string        `json:"ip"`
}

type SellerProfile struct {
	SellerID        int    `json:"seller_id"`
	SellerUUID      string `json:"seller_uuid"`
	Role            string `json:"role"`
	Name            string `json:"name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	MerchantID      int    `json:"merchant_id"`
	ProfileImageURL string `json:"profile_image_url"`
	RegisterDate    string `json:"register_date"`
	LastLoginDate   string `json:"last_login_date"`
}

type RefreshTokenData struct {
	RefreshTokenID  int            `json:"refresh_token_id"`
	SellerID        int            `json:"seller_id"`
	SellerUUID      string         `json:"seller_uuid"`
	Name            string         `json:"name"`
	LastName        string         `json:"last_name"`
	Email           string         `json:"email"`
	Password        string         `json:"password"`
	ProfileImageURL sql.NullString `json:"profile_image_url"`
	Token           string         `json:"token"`
	TokenType       string         `json:"token_type"`
	ExpiredDate     string         `json:"expired_date"`
}

type SellerAccountDataLogin struct {
	SellerID        int    `json:"seller_id"`
	SellerUUID      string `json:"seller_uuid"`
	Role            string `json:"role"`
	Name            string `json:"name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	MerchantID      int64  `json:"merchant_id"`
	ProfileImageURL string `json:"profile_image_url"`
	RefreshToken    string `json:"refresh_token"`
	AccessToken     string `json:"access_token"`
	CreateDate      string `json:"create_date"`
	ExpiredDate     string `json:"expired_date"`
}
