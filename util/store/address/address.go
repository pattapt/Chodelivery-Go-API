package UtilStore_Address

import (
	"fmt"
	"gobasic/database"
	StoreAddressStruct "gobasic/struct/store/address"
)

func GetAddressByToken(AddressToken string, OwnerId int) (StoreAddressStruct.Destination, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return StoreAddressStruct.Destination{}, err
	}
	defer db.Close()

	var Destination StoreAddressStruct.Destination
	query := `
		SELECT d.DestinationId, d.DestinationToken, d.Name, d.PhoneNumber, d.Status, d.CreateDate, d.UpdateDate,
		d.Address, d.Street, d.Building, dt.id as DistrictId, dt.name_th as district, a.name_th as amphure, p.name_th as province, 
		dt.zip_code as zipcode, d.Note
		FROM Destination d
		LEFT JOIN districts dt ON d.distric = dt.id
		LEFT JOIN amphures a ON dt.amphure_id = a.id
		LEFT JOIN provinces p ON a.province_id = p.id
		WHERE d.DestinationToken = ? AND d.AccountId = ? AND d.Status = 'show'
    `
	err = db.QueryRow(query, AddressToken, OwnerId).Scan(
		&Destination.DestinationId,
		&Destination.DestinationToken,
		&Destination.Name,
		&Destination.PhoneNumber,
		&Destination.Status,
		&Destination.CreateDate,
		&Destination.UpdateDate,
		&Destination.Address.Address,
		&Destination.Address.Street,
		&Destination.Address.Building,
		&Destination.Address.DistrictId,
		&Destination.Address.District,
		&Destination.Address.Amphure,
		&Destination.Address.Province,
		&Destination.Address.ZipCode,
		&Destination.Note,
	)
	if err != nil {
		fmt.Print(err.Error())
		return StoreAddressStruct.Destination{}, err
	}

	return Destination, nil
}

func GetAddressById(AddressId int, OwnerId int) (StoreAddressStruct.Destination, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return StoreAddressStruct.Destination{}, err
	}
	defer db.Close()

	var Destination StoreAddressStruct.Destination
	query := `
		SELECT d.DestinationId, d.DestinationToken, d.Name, d.PhoneNumber, d.Status, d.CreateDate, d.UpdateDate,
		d.Address, d.Street, d.Building, dt.id as DistrictId, dt.name_th as district, a.name_th as amphure, p.name_th as province, 
		dt.zip_code as zipcode, d.Note
		FROM Destination d
		LEFT JOIN districts dt ON d.distric = dt.id
		LEFT JOIN amphures a ON dt.amphure_id = a.id
		LEFT JOIN provinces p ON a.province_id = p.id
		WHERE d.DestinationId = ? AND d.AccountId = ? AND d.Status = 'show'
    `
	err = db.QueryRow(query, AddressId, OwnerId).Scan(
		&Destination.DestinationId,
		&Destination.DestinationToken,
		&Destination.Name,
		&Destination.PhoneNumber,
		&Destination.Status,
		&Destination.CreateDate,
		&Destination.UpdateDate,
		&Destination.Address.Address,
		&Destination.Address.Street,
		&Destination.Address.Building,
		&Destination.Address.DistrictId,
		&Destination.Address.District,
		&Destination.Address.Amphure,
		&Destination.Address.Province,
		&Destination.Address.ZipCode,
		&Destination.Note,
	)
	if err != nil {
		fmt.Print(err.Error())
		return StoreAddressStruct.Destination{}, err
	}

	return Destination, nil
}
