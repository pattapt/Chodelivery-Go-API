package MerchantV1Store

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	config "gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
)

func GetSummary(c *gin.Context) {

	db, err := database.ConnectDB()
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer db.Close()

	merchantData, exists := c.Get("merchantData")
	if !exists {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	mcData, ok := merchantData.(MerchantMerchantStruct.MerchantDetail)
	if !ok {
		errorResponse := APIStructureMain.ErrorRespondAuth("ข้อมูลไม่ถูกต้อง", "Invalid data")
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	currentTime := time.Now()
	minus24Hours := currentTime.Add(-24 * time.Hour)
	minus7Days := currentTime.Add(-7 * 24 * time.Hour)
	Minus24Hour := minus24Hours.Format("2006-01-02 15:04:05")
	Minus7Day := minus7Days.Format("2006-01-02 15:04:05")

	var Summary MerchantMerchantStruct.MerchantSummary
	query := `
		SELECT COUNT(1) AS TotalOrder, SUM(TotalPay) AS TotalAmount FROM Transaction 
		WHERE MerchantId = ? AND Status != 'failed' AND CreateDate BETWEEN ? AND ?
	`
	err = db.QueryRow(query, mcData.MerchantId, Minus24Hour, config.GetCurrentDateTime()).Scan(
		&Summary.TotalOrder,
		&Summary.TotalAmount,
	)
	query = `
		SELECT COUNT(1) AS LastWeekOrder, SUM(TotalPay) AS LastWeekSales FROM Transaction 
		WHERE MerchantId = ? AND Status != 'failed' AND CreateDate BETWEEN ? AND ?
	`
	err = db.QueryRow(query, mcData.MerchantId, Minus7Day, config.GetCurrentDateTime()).Scan(
		&Summary.LastWeekOrder,
		&Summary.LastWeekSales,
	)

	// Get Profile
	respond := APIStructureMain.Response{
		StatusCode: 200,
		Status:     "S",
		Message:    "Successfully Get Data",
		Data:       Summary,
	}
	c.JSON(http.StatusOK, respond)

}
