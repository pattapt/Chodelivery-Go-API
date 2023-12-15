package StoreV1Merchant

import (
	"net/http"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"

	StoreMerchantStruct "gobasic/struct/store/merchant"
)

func GetMerchantInfo(c *gin.Context) {
	MerchantId, exists := c.Get("merchantId")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer db.Close()

	var merchant StoreMerchantStruct.MerchantDetail
	query := `SELECT m.MerchantId, m.MerchantUUID, m.Name, m.Description,
				CASE WHEN m.Status = 'open' THEN true ELSE false END AS Open,
				CASE WHEN m.Visible = 'visible' THEN true ELSE false END AS Visible,
				m.ImageUrl, m.Address, m.Street, m.Building, 
				d.name_th as district, a.name_th as amphure, p.name_th as province, d.zip_code as zipcode, 
				m.CreateDate, m.UpdateDate
				FROM Merchant m
				LEFT JOIN districts d ON m.distric = d.id
				LEFT JOIN amphures a ON d.amphure_id = a.id
				LEFT JOIN provinces p ON a.province_id = p.id
				WHERE MerchantId = ?`
	err = db.QueryRow(query, MerchantId).Scan(
		&merchant.MerchantId,
		&merchant.MerchantUUID,
		&merchant.Name,
		&merchant.Description,
		&merchant.Open,
		&merchant.Visible,
		&merchant.ImageUrl,
		&merchant.Address.Address,
		&merchant.Address.Street,
		&merchant.Address.Building,
		&merchant.Address.District,
		&merchant.Address.Amphure,
		&merchant.Address.Province,
		&merchant.Address.ZipCode,
		&merchant.CreateDate,
		&merchant.UpdateDate,
	)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	merchant.ImageUrl = config.GetMainAPIURL() + merchant.ImageUrl

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully retrieved data",
		Data:       merchant,
	}
	c.JSON(http.StatusOK, respond)
}
