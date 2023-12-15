package APIStructureMain

func DEBUG(msg string) Response {
	return ErrorRespondAuth("ไม่สามารถทำรายการได้", msg)
}

func GeneralErrorMSG() Response {
	return ErrorRespondAuth("ไม่สามารถทำรายการได้", "ระบบไม่สามารถทำรายการได้ในขณะนี้ โปรดลองใหม่ภายหลัง")
}

func AccessTokenNotValid() Response {
	return ErrorRespondAuth("ไม่สามารถทำรายการได้", "ข้อมูล Access Token ของท่านหมดอายุแล้ว กรุณาทำการเข้าสู่ระบบใหม่")
}

func ErrorRespondAuth(args ...string) Response {
	var title, desc string

	if len(args) > 0 {
		title = args[0]
	} else {
		title = "ไม่สามารถทำรายการได้ในขณะนี้"
	}

	if len(args) > 1 {
		desc = args[1]
	} else {
		desc = "ข้อมูลของคุณไม่เพียงพอ กรุณาแก้ไขข้อมูลของคุณ และทำรายการใหม่อีกครั้ง"
	}

	errorData := ErrorData{
		Title:       title,
		Description: desc,
	}

	errorResponse := Response{
		StatusCode: 402,
		Status:     "E",
		Message:    "Invalid Data",
		Data:       errorData,
	}

	return errorResponse
}

type ErrorData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Response struct {
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Message    string      `json:"msg"`
	Data       interface{} `json:"data"`
}

type RegisterAccountStoreSuccess struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AccountDataLogin struct {
	AccountID       int    `json:"account_id"`
	AccountUUID     string `json:"account_uuid"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	ProfileImageURL string `json:"profile_image_url"`
	RefreshToken    string `json:"refresh_token"`
	AccessToken     string `json:"access_token"`
	CreateDate      string `json:"create_date"`
	ExpiredDate     string `json:"expired_date"`
}

type AccountProfile struct {
	AccountID       int    `json:"account_id"`
	AccountUUID     string `json:"account_uuid"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	ProfileImageURL string `json:"profile_image_url"`
	RegisterDate    string `json:"RegisterDate"`
	LastLoginDate   string `json:"LastLoginDate"`
}
