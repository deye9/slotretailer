package service

import (
	"database/sql"
)

// LocalStore returns a pointer to the local Store object
var LocalStore Store

// GetStore returns the details of the store
func GetStore() (store Store, err error) {
	var rows *sql.Rows
	if rows, err = Get(`select * from store;`); err != nil {
		CheckError("Error getting store data.", err, false)
		return Store{}, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&store.ID, &store.Name, &store.Address, &store.Phone, &store.City, &store.Email, &store.OrdersAPI, &store.ProductsAPI, &store.CustomersAPI, &store.CreditCardAPI, &store.SyncInterval, &store.SapKey, &store.CreatedBy, &store.CreatedAt, &store.UpdatedAt, &store.DeletedAt, &store.LogRotation, &store.TransfersAPI, &store.AllowVAT, &store.WarehousesAPI, &store.PricelistAPI, &store.ProductPriceList, &store.CashAccountAPI, &store.StoreCashAccount, &store.BankTransferAPI, &store.ChequesAPI); err != nil {
			CheckError("Error Scanning store.", err, false)
		} else {
			LocalStore = store
		}
	}

	return
}

// SaveStore creates / updates the store in the database
func SaveStore(store map[string]interface{}) (id int64, err error) {

	if store["id"] == nil {
		if result, err := MaptoInsert(store, "store"); err == nil {
			if id, err = Insert(result); err != nil {
				CheckError("Error inserting the store data.", err, false)
			}
		}
		return
	}

	if result, err := MaptoUpdate(store, "store", "id"); err == nil {
		id = int64(store["id"].(float64))
		if err = Modify(result); err != nil {
			CheckError("Error updating the store data.", err, false)
		}
	}
	return
}

// GetStores returns an array of Stores
func GetStores() (stores []Stores, err error) {
	var rows *sql.Rows

	if rows, err = Get(`select * from stores order by name;`); err != nil {
		CheckError("Error getting available Stores.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		store := Stores{}
		if err = rows.Scan(&store.ID, &store.Name, &store.Code); err != nil {
			CheckError("Error Scanning stores.", err, false)
		} else {
			stores = append(stores, store)
		}
	}

	return
}
