package service

import (
	"database/sql"
	"fmt"
	"strings"
)

// 1. 7 days return policy

// ProductDetails returns a instance belonging to the Product id passed in
func ProductDetails(id int) (product Products, err error) {
	var rows *sql.Rows
	if rows, err = Get(fmt.Sprintf(`select * from products where id = %d;`, id)); err != nil {
		CheckError("Error getting Product data.", err, false)
		return Products{}, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&product.ID, &product.ItemCode, &product.ItemName, &product.CodeBars, &product.OnHand, &product.MinLevel, &product.Warehouse, &product.SerialNumber, &product.Manufacturer, &product.Price, &product.Vat, &product.ItemID); err != nil {
			CheckError("Error Scanning Product.", err, false)
		}
	}

	return
}

// GetProducts returns an array of Products
func GetProducts() (products []Products, err error) {
	var rows *sql.Rows

	if rows, err = Get(`select * from products where warehouse = "` + LocalStore.SapKey + `" order by ID;`); err != nil {
		CheckError("Error getting Products.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		product := Products{}
		if err = rows.Scan(&product.ID, &product.ItemCode, &product.ItemName, &product.CodeBars, &product.OnHand, &product.MinLevel, &product.Warehouse, &product.SerialNumber, &product.Manufacturer, &product.Price, &product.Vat, &product.ItemID); err != nil {
			CheckError("Error Scanning Products.", err, false)
		} else {
			products = append(products, product)
		}
	}

	return
}

// GetStoreProducts returns an array of Products belonging to the specified Store
func GetStoreProducts(id int) (products []Products, err error) {
	var rows *sql.Rows

	if rows, err = Get(fmt.Sprintf(`SELECT p.* FROM stores s inner join products p on s.name = p.warehouse where s.id = %d order by p.itemcode;`, id)); err != nil {
		CheckError("Error getting Store Products.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		product := Products{}
		if err = rows.Scan(&product.ID, &product.ItemCode, &product.ItemName, &product.CodeBars, &product.OnHand, &product.MinLevel, &product.Warehouse, &product.SerialNumber, &product.Manufacturer, &product.Price, &product.Vat, &product.ItemID); err != nil {
			CheckError("Error Scanning Store Products.", err, false)
		} else {
			products = append(products, product)
		}
	}

	return
}

// GetTransfer returns a instance belonging to the Transfer id passed in
func GetTransfer(id int) (transfer Transfers, err error) {
	var rows *sql.Rows
	if rows, err = Get(fmt.Sprintf(`select t.*, i.id, i.transferid, i.itemcode, i.itemname, p.onhand, i.quantity, p.serialnumbers from transfers t inner join transfereditems i on t.id = i.transferid inner join products p on p.itemcode = i.itemcode where t.id = %d and t.deleted_at is null;`, id)); err != nil {
		CheckError("Error getting Transfer Request data.", err, false)
		return Transfers{}, err
	}
	defer rows.Close()

	var items []Transfereditems
	for rows.Next() {
		item := Transfereditems{}

		if err = rows.Scan(&transfer.ID, &transfer.FromWhs, &transfer.ToWhs, &transfer.Comment, &transfer.Canceled, &transfer.Synced,
			&transfer.Status, &transfer.CreatedBy, &transfer.CreatedAt, &transfer.UpdatedAt, &transfer.DeletedAt, &transfer.DocEntry, &transfer.DocNum,
			&item.ID, &item.TransferID, &item.ItemCode, &item.ItemName, &item.OnHand, &item.Quantity, &item.SerialNumber); err != nil {
			CheckError("Error Scanning Transfer Request.", err, false)
		} else {
			items = append(items, item)
		}
	}

	transfer.Items = items
	return
}

// GetTransfers returns an array of Transfers
func GetTransfers() (transfers []Transfers, err error) {
	var rows *sql.Rows
	if rows, err = Get(`select * from transfers where deleted_at is null order by created_at desc;`); err != nil {
		CheckError("Error getting Transfers.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		transfer := Transfers{}
		if err = rows.Scan(&transfer.ID, &transfer.FromWhs, &transfer.ToWhs, &transfer.Comment, &transfer.Canceled, &transfer.Synced, &transfer.Status, &transfer.CreatedBy, &transfer.CreatedAt, &transfer.UpdatedAt, &transfer.DeletedAt, &transfer.DocEntry, &transfer.DocNum); err != nil {
			CheckError("Error Scanning Transfer.", err, false)
		} else {
			transfers = append(transfers, transfer)
		}
	}

	return
}

// NewTransfer creates a new Transfer in the database
func NewTransfer(transfer map[string]interface{}) (err error) {
	master := ""
	detail := ""

	// Get a new reference to the transfered items and remove it from the map.
	itemsOrdered := transfer["items"]
	delete(transfer, "items")

	// Get the master insert query
	if master, err = MaptoInsert(transfer, "transfers"); err != nil {
		CheckError("Error Mapping the Transfer to SQL.", err, false)
		return err
	}

	master += "SET @last_id := (SELECT LAST_INSERT_ID());"

	// Get the details insert query
	for _, _value := range itemsOrdered.([]interface{}) {
		if detail, err = MaptoInsert(_value.(map[string]interface{}), "transfereditems"); err != nil {
			CheckError("Error Mapping the Transfered Items to SQL.", err, false)
			return err
		}

		// Build out the needed queries
		master += strings.Replace(fmt.Sprintf("%v", detail), `""`, "@last_id", -1)
	}

	// Save the Order and Reduce inventory
	if err = Modify(master); err != nil {
		CheckError("Error creating the Transfer.", err, false)
		return err
	}
	return
}

// UpdateTransfer updates the Transfer details in the database
func UpdateTransfer(transfer map[string]interface{}) (id int, err error) {
	// Get a new reference to the transfered items and remove it from the map.
	var result, detail string
	itemsOrdered := transfer["items"]
	delete(transfer, "items")

	if result, err = MaptoUpdate(transfer, "transfers", "id"); err != nil {
		CheckError("Error generating the Update Transfer SQL.", err, false)
		return 0, err
	}

	for _, _value := range itemsOrdered.([]interface{}) {
		item := _value.(map[string]interface{})

		if item["id"] == 0.00 {
			// Generate the Insert Query
			if detail, err = MaptoInsert(item, "transfereditems"); err != nil {
				CheckError("Error Mapping the Transfered Items to an INSERT SQL.", err, false)
				return 0, err
			}
		} else {
			// Generate the Update Query
			if detail, err = MaptoUpdate(item, "transfereditems", "id"); err != nil {
				CheckError("Error Mapping the Transfered Items to a UPDATE SQL.", err, false)
				return 0, err
			}
		}
		// Build out the needed queries
		result += strings.Replace(fmt.Sprintf("%v", detail), `updated_at = CURRENT_TIMESTAMP,`, "", -1)
	}

	if err = Modify(result); err != nil {
		CheckError("Error updating the Transfer.", err, false)
	}
	return
}

// RemoveTransfer deletes a Transfer from the database
func RemoveTransfer(id int) (err error) {
	if err = Modify(fmt.Sprintf(`update transfers set deleted_at = CURRENT_TIMESTAMP where id = %d;`, id)); err != nil {
		CheckError("Error removing Transfer.", err, false)
		return err
	}

	return nil
}
