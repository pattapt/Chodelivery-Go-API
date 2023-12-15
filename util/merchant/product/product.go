package UtilMerchant_Product

import (
	"gobasic/config"
	"gobasic/database"
	MerchantProductStruct "gobasic/struct/merchant/product"
)

func GetProductDetail(ProductId string, MerchantId int) (MerchantProductStruct.Product, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return MerchantProductStruct.Product{}, err
	}
	defer db.Close()
	query := `
		SELECT p.ProductId, p.ProductToken, p.Barcode, p.CategoryId, 
		p.Name, p.Description, p.ImageUrl, p.Price, p.Cost, p.Quantity, 
		p.Status,
		CASE WHEN p.Invisible = 'visible' THEN true ELSE false END AS visible,
		p.CreateDate, p.UpdateDate
		FROM Product p WHERE p.ProductToken = ? AND p.MerchantId = ? 
	`
	var pd MerchantProductStruct.Product
	err = db.QueryRow(query, ProductId, MerchantId).Scan(
		&pd.ProductId,
		&pd.ProductToken,
		&pd.Barcode,
		&pd.CategoryId,
		&pd.Name,
		&pd.Description,
		&pd.ImageUrl,
		&pd.Price,
		&pd.Cost,
		&pd.StockQuantity,
		&pd.Status,
		&pd.Visible,
		&pd.CreateDate,
		&pd.UpdateDate,
	)
	if err != nil {
		return MerchantProductStruct.Product{}, err
	}
	pd.ImageUrl = config.GetMainAPIURL() + pd.ImageUrl

	return pd, nil
}
