package StoreAccountStruct

import "database/sql"

type RegisterAccountPost struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type LoginAccountPost struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type InvokeAccessToken struct {
	RefreshToken string `json:"refreshToken"`
	UUID         string `json:"uuid"`
}

type AccountData struct {
	AccountID       int            `json:"account_id"`
	AccountUUID     string         `json:"account_uuid"`
	Username        string         `json:"username"`
	Email           string         `json:"email"`
	Password        string         `json:"password"`
	ProfileImageURL sql.NullString `json:"profile_image_url"`
	RegisterDate    string         `json:"register_date"`
	LastLoginDate   string         `json:"last_login_date"`
}

type RefreshTokenData struct {
	RefreshTokenID  int            `json:"refresh_token_id"`
	AccountID       int            `json:"account_id"`
	AccountUUID     string         `json:"account_uuid"`
	Username        string         `json:"username"`
	Email           string         `json:"email"`
	Password        string         `json:"password"`
	ProfileImageURL sql.NullString `json:"profile_image_url"`
	Token           string         `json:"token"`
	TokenType       string         `json:"token_type"`
	ExpiredDate     string         `json:"expired_date"`
}

type AccessTokenData struct {
	AccessTokenID  int    `json:"access_token_id"`
	AccountID      int    `json:"account_id"`
	TokenType      string `json:"token_type"`
	AccountType    string `json:"account_type"`
	RefreshTokenID int    `json:"refresh_token_id"`
	Token          string `json:"token"`
	IssueDate      string `json:"issue_date"`
	ExpiredDate    string `json:"expired_date"`
	IP             string `json:"ip"`
}
