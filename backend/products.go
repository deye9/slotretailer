package service

import (
	"database/sql"
	"fmt"
)

// ProductDetails returns a instance belonging to the Product id passed in
func ProductDetails(id int) (product Products, err error) {
	var rows *sql.Rows
	if rows, err = Get(fmt.Sprintf(`select * from products where id = %d;`, id)); err != nil {
		CheckError("Error getting Product data.", err, false)
		return Products{}, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&product.ID, &product.ItemCode, &product.ItemName, &product.CodeBars, &product.OnHand, &product.MinLevel, &product.Warehouse, &product.SerialNumber, &product.Manufacturer, &product.Price, &product.Vat); err != nil {
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
		if err = rows.Scan(&product.ID, &product.ItemCode, &product.ItemName, &product.CodeBars, &product.OnHand, &product.MinLevel, &product.Warehouse, &product.SerialNumber, &product.Manufacturer, &product.Price, &product.Vat); err != nil {
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
		if err = rows.Scan(&product.ID, &product.ItemCode, &product.ItemName, &product.CodeBars, &product.OnHand, &product.MinLevel, &product.Warehouse, &product.SerialNumber, &product.Manufacturer, &product.Price, &product.Vat); err != nil {
			CheckError("Error Scanning Store Products.", err, false)
		} else {
			products = append(products, product)
		}
	}

	return
}
