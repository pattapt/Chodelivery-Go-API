package MerchantV1Store

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
)

func GetAllMerchant(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	pageStr := c.DefaultQuery("page", "1")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 50
	}

	Start := (page - 1) * limit

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

	// Count total merchant
	var totalMerchant int
	err = db.QueryRow(`SELECT COUNT(*) FROM Merchant WHERE OwnerSellerId = ?`, sellerId).Scan(&totalMerchant)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where no rows are found
			errorResponse := APIStructureMain.ErrorRespondAuth("Merchant not found", "No merchant found for the given seller ID")
			c.JSON(http.StatusNotFound, errorResponse)
			return
		}
		// Handle other errors
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// Calculate total page
	totalPage := totalMerchant / limit
	if totalMerchant%limit != 0 {
		totalPage++
	}

	query := `SELECT m.MerchantId, m.MerchantUUID, m.Name, m.Description,
				CASE WHEN m.Status = 'open' THEN true ELSE false END AS Open,
				CASE WHEN m.Visible = 'visible' THEN true ELSE false END AS Visible,
				m.ImageUrl, m.Address, m.Street, m.Building, 
				d.id AS DistrictId, d.name_th as district, a.name_th as amphure, p.name_th as province, d.zip_code as zipcode, 
				m.CreateDate, m.UpdateDate
				FROM Merchant m
				LEFT JOIN districts d ON m.distric = d.id
				LEFT JOIN amphures a ON d.amphure_id = a.id
				LEFT JOIN provinces p ON a.province_id = p.id
				WHERE m.OwnerSellerId = ?
				ORDER BY m.MerchantId DESC
				LIMIT ?, ?`
	rows, err := db.Query(query, sellerId, Start, limit)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer rows.Close()

	var merchantList []MerchantMerchantStruct.MerchantDetail

	for rows.Next() {
		var merchant MerchantMerchantStruct.MerchantDetail

		err := rows.Scan(
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
			&merchant.Address.DistrictId,
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
		merchantList = append(merchantList, merchant)
	}

	if len(merchantList) == 0 {
		// Return an empty array as the response
		respond := APIStructureMain.Response{
			StatusCode: http.StatusOK,
			Status:     "S",
			Message:    "No merchant found",
			Data: MerchantMerchantStruct.MerchantList{
				Merchants: []MerchantMerchantStruct.MerchantDetail{},
			},
		}
		c.JSON(http.StatusOK, respond)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully get merchant",
		Data: MerchantMerchantStruct.MerchantList{
			PageId:    page,
			TotalPage: totalMerchant,
			TotalItem: totalMerchant,
			Merchants: merchantList,
		},
	}
	c.JSON(http.StatusOK, respond)
}
