package MerchantV1Store

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
)

func GetMerchantDetail(c *gin.Context) {
	MerchantUUID := c.Param("MerchantUUID")
	if MerchantUUID == "" {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer db.Close()

	sellerId, exists := c.Get("sellerId")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var merchant MerchantMerchantStruct.MerchantDetail
	query := `SELECT m.MerchantId, m.MerchantUUID, m.Name, m.Description, m.PromptpayPhone,
				CASE WHEN m.Status = 'open' THEN true ELSE false END AS Open,
				CASE WHEN m.Visible = 'visible' THEN true ELSE false END AS Visible,
				m.ImageUrl, m.Address, m.Street, m.Building, 
				d.id AS DistrictId, d.name_th as district, a.name_th as amphure, p.name_th as province, d.zip_code as zipcode, 
				m.CreateDate, m.UpdateDate
				FROM Merchant m
				LEFT JOIN districts d ON m.distric = d.id
				LEFT JOIN amphures a ON d.amphure_id = a.id
				LEFT JOIN provinces p ON a.province_id = p.id
				WHERE m.MerchantUUID = ? AND m.OwnerSellerId = ?`
	err = db.QueryRow(query, MerchantUUID, sellerId).Scan(
		&merchant.MerchantId,
		&merchant.MerchantUUID,
		&merchant.Name,
		&merchant.Description,
		&merchant.PromptPayPhone,
		&merchant.Open,
		&merchant.Visible,
		&merchant.ImageUrl,
		&merchant.Address.Address,
		&merchant.Address.Street,
		&merchant.Address.Building,
		&merchant.Address.DistrictId,
		&merchant.Address.District,
		&merchant.Address.Amphure,
		&merchant.Address.Province,
		&merchant.Address.ZipCode,
		&merchant.CreateDate,
		&merchant.UpdateDate,
	)
	if err != nil {
		fmt.Print(err)
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	merchant.ImageUrl = config.GetMainAPIURL() + merchant.ImageUrl

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully get merchant",
		Data:       merchant,
	}
	c.JSON(http.StatusOK, respond)
}
