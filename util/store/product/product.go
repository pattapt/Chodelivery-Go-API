package UtilStore_Product

import (
	"fmt"
	"gobasic/config"
	"gobasic/database"
	StoreCategoryStruct "gobasic/struct/store/merchant"
)

func GetProductById(ProductId int) (StoreCategoryStruct.Product, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return StoreCategoryStruct.Product{}, err
	}
	defer db.Close()

	var product StoreCategoryStruct.Product
	query := `
        SELECT 
            ProductId,
            ProductToken,
            Barcode,
            CategoryId AS ProductCategoryId,
            Name AS ProductName,
            Description AS ProductDescription,
            ImageUrl AS ProductImageUrl,
            Price AS ProductPrice,
            Status AS ProductStatus,
            CreateDate AS ProductCreateDate,
            UpdateDate AS ProductUpdateDate
        FROM Product
        WHERE ProductId = ?
    `
	err = db.QueryRow(query, ProductId).Scan(
		&product.ProductId,
		&product.ProductToken,
		&product.Barcode,
		&product.ProductCategoryId,
		&product.Name,
		&product.Description,
		&product.ImageUrl,
		&product.Price,
		&product.Status,
		&product.CreateDate,
		&product.UpdateDate,
	)
	if err != nil {
		fmt.Print(err.Error())
		return StoreCategoryStruct.Product{}, err
	}

	return product, nil
}

func GetMerchantInfoById(MerchantId int) (StoreCategoryStruct.MerchantDetailV2, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return StoreCategoryStruct.MerchantDetailV2{}, err
	}
	defer db.Close()
	var merchant StoreCategoryStruct.MerchantDetailV2
	query := `SELECT m.MerchantId, m.MerchantUUID, m.OwnerSellerId, m.Name, m.Description,
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
		&merchant.OwnerID,
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
		return StoreCategoryStruct.MerchantDetailV2{}, err
	}

	merchant.ImageUrl = config.GetMainAPIURL() + merchant.ImageUrl

	return merchant, nil
}
