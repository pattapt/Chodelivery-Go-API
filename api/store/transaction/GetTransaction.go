package StoreV1Transaction

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	StoreTransactionStruct "gobasic/struct/store/transaction"
)

func GetTransaction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

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

	query := `SELECT t.OrderId, t.OrderToken, t.PaymentMethod, t.TotalPay, t.Status, 
				t.Note, t.ChatId, t.CreateDate, t.UpdateDate,
				m.MerchantId, m.MerchantUUID, m.OwnerSellerId, m.Name, m.Description,
				CASE WHEN m.Status = 'open' THEN true ELSE false END AS Open,
				CASE WHEN m.Visible = 'visible' THEN true ELSE false END AS Visible,
				m.ImageUrl, m.Address, m.Street, m.Building, 
				d.name_th as district, a.name_th as amphure, p.name_th as province, d.zip_code as zipcode
				FROM Transaction t
				LEFT JOIN Merchant m ON t.MerchantId = m.MerchantId
				LEFT JOIN districts d ON m.distric = d.id
				LEFT JOIN amphures a ON d.amphure_id = a.id
				LEFT JOIN provinces p ON a.province_id = p.id
				WHERE t.AccountId = ?
				ORDER BY t.OrderId DESC
				LIMIT ?, ?
				`
	rows, err := db.Query(query, userID, Start, limit)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer rows.Close()

	var transactionList []StoreTransactionStruct.TransactionList

	for rows.Next() {
		var tran StoreTransactionStruct.TransactionList

		err := rows.Scan(
			&tran.OrderId,
			&tran.OrderToken,
			&tran.PaymentMethod,
			&tran.TotalPay,
			&tran.Status,
			&tran.Note,
			&tran.ChatId,
			&tran.CreateDate,
			&tran.UpdateDate,
			&tran.Merchant.MerchantId,
			&tran.Merchant.MerchantUUID,
			&tran.Merchant.OwnerSellerId,
			&tran.Merchant.Name,
			&tran.Merchant.Description,
			&tran.Merchant.Open,
			&tran.Merchant.Visible,
			&tran.Merchant.ImageUrl,
			&tran.Merchant.Address,
			&tran.Merchant.Street,
			&tran.Merchant.Building,
			&tran.Merchant.District,
			&tran.Merchant.Amphure,
			&tran.Merchant.Province,
			&tran.Merchant.ZipCode,
		)
		tran.Merchant.ImageUrl = config.GetMainAPIURL() + tran.Merchant.ImageUrl

		if err != nil {
			errorResponse := APIStructureMain.GeneralErrorMSG()
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		transactionList = append(transactionList, tran)
	}

	if len(transactionList) == 0 {
		// Return an empty array as the response
		respond := APIStructureMain.Response{
			StatusCode: http.StatusOK,
			Status:     "S",
			Message:    "No transaction found",
			Data:       []StoreTransactionStruct.TransactionList{},
		}
		c.JSON(http.StatusOK, respond)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully get transaction",
		Data:       transactionList,
	}
	c.JSON(http.StatusOK, respond)
}
