package MerchantV1Transaction

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gobasic/config"
	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantMerchantStruct "gobasic/struct/merchant/merchant"
	MerchantTransactionStruct "gobasic/struct/merchant/transaction"
)

func GetTransaction(c *gin.Context) {
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
				cs.AccountId, cs.AccountUUID, cs.Username, cs.ProfileImageURL,
				ds.DestinationId, ds.DestinationToken, ds.Name AS DestinationName, ds.Phonenumber,
				ds.Address, ds.Street, ds.Building, ds.Note,
				d.name_th as district, a.name_th as amphure, p.name_th as province, d.zip_code as zipcode
				FROM Transaction t
				LEFT JOIN Merchant m ON t.MerchantId = m.MerchantId
				LEFT JOIN Account cs ON cs.AccountId = t.AccountId
				LEFT JOIN Destination ds ON t.DestinationId = ds.DestinationId
				LEFT JOIN districts d ON ds.distric = d.id
				LEFT JOIN amphures a ON d.amphure_id = a.id
				LEFT JOIN provinces p ON a.province_id = p.id
				WHERE t.MerchantId = ?
				ORDER BY t.OrderId DESC
				LIMIT ?, ?;
				`
	rows, err := db.Query(query, mcData.MerchantId, Start, limit)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer rows.Close()

	var transactionList []MerchantTransactionStruct.TransactionList

	for rows.Next() {
		var tran MerchantTransactionStruct.TransactionList

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
			&tran.Customer.AccountId,
			&tran.Customer.AccountUUID,
			&tran.Customer.Username,
			&tran.Customer.ProfileImageURL,
			&tran.Destination.DestinationId,
			&tran.Destination.DestinationToken,
			&tran.Destination.DestinationName,
			&tran.Destination.PhoneNumber,
			&tran.Destination.Address,
			&tran.Destination.Street,
			&tran.Destination.Building,
			&tran.Destination.Note,
			&tran.Destination.District,
			&tran.Destination.Amphure,
			&tran.Destination.Province,
			&tran.Destination.ZipCode,
		)
		tran.Customer.ProfileImageURL = config.GetMainAPIURL() + tran.Customer.ProfileImageURL

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
			Data:       []MerchantTransactionStruct.TransactionList{},
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
