package MerchantV1Util

import (
	"net/http"

	"github.com/gin-gonic/gin"

	database "gobasic/database"
	APIStructureMain "gobasic/struct"
	MerchantProductStruct "gobasic/struct/merchant/product"
	MerchantUtilStruct "gobasic/struct/merchant/util"
)

func GetDistrict(c *gin.Context) {
	Zipcode := c.DefaultQuery("zipcode", "1")

	db, err := database.ConnectDB()
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer db.Close()

	query := `SELECT d.id AS DistrictId, d.name_th AS DistrictNameTh, d.name_en AS DistrictNameEn,
				a.id AS AmphureId, a.name_th AS AmphureNameTh, a.name_en AS AmphureNameEn, 
				p.id AS ProvinceId, p.name_th AS ProvinceNameTh, p.name_en AS ProvinceNameEn,
				d.zip_code AS ZipCode
				FROM districts d 
				INNER JOIN amphures a ON d.amphure_id = a.id 
				INNER JOIN provinces p ON a.province_id = p.id
				WHERE d.zip_code = ?`
	rows, err := db.Query(query, Zipcode)
	if err != nil {
		errorResponse := APIStructureMain.GeneralErrorMSG()
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	defer rows.Close()

	var Provinces []MerchantUtilStruct.THProvince

	for rows.Next() {
		var data MerchantUtilStruct.THProvince

		err := rows.Scan(
			&data.District.DistrictId,
			&data.District.NameTh,
			&data.District.NameEn,
			&data.Amphure.AmphureId,
			&data.Amphure.NameTh,
			&data.Amphure.NameEn,
			&data.Province.ProvinceId,
			&data.Province.NameTh,
			&data.Province.NameEn,
			&data.ZipCode,
		)

		if err != nil {
			// fmt.Print(err)
			errorResponse := APIStructureMain.GeneralErrorMSG()
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		Provinces = append(Provinces, data)
	}

	if len(Provinces) == 0 {
		// Return an empty array as the response
		respond := APIStructureMain.Response{
			StatusCode: http.StatusOK,
			Status:     "S",
			Message:    "No Data found",
			Data:       []MerchantProductStruct.Product{},
		}
		c.JSON(http.StatusOK, respond)
		return
	}

	respond := APIStructureMain.Response{
		StatusCode: http.StatusOK,
		Status:     "S",
		Message:    "Successfully get Data",
		Data:       Provinces,
	}
	c.JSON(http.StatusOK, respond)
}
