package service

import (
	"database/sql"
	"fmt"
)

// GetCustomer returns a instance belonging to the Customers id passed in
func GetCustomer(id int) (customer Customers, err error) {
	var rows *sql.Rows
	if rows, err = Get(fmt.Sprintf(`select * from customers where id = %d;`, id)); err != nil {
		CheckError("Error getting Customer data.", err, false)
		return Customers{}, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&customer.ID, &customer.CardCode, &customer.CardName, &customer.Address, &customer.Phone, &customer.Phone1, &customer.City, &customer.Email, &customer.Synced, &customer.CreatedBy, &customer.CreatedAt, &customer.UpdatedAt, &customer.DeletedAt); err != nil {
			CheckError("Error Scanning Customer.", err, false)
		}
	}

	return
}

// GetCustomers returns an array of Customers
func GetCustomers() (customers []Customers, err error) {
	var rows *sql.Rows
	if rows, err = Get(`select * from customers where deleted_at is null;`); err != nil {
		CheckError("Error getting Customers.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		customer := Customers{}
		if err = rows.Scan(&customer.ID, &customer.CardCode, &customer.CardName, &customer.Address, &customer.Phone, &customer.Phone1, &customer.City, &customer.Email, &customer.Synced, &customer.CreatedBy, &customer.CreatedAt, &customer.UpdatedAt, &customer.DeletedAt); err != nil {
			CheckError("Error Scanning Customers.", err, false)
		} else {
			customers = append(customers, customer)
		}
	}

	return
}

// GetCustomerbyPhone returns a instance belonging to the Customers id passed in
func GetCustomerbyPhone(details string) (customer Customers, err error) {
	var rows *sql.Rows
	if rows, err = Get(fmt.Sprintf(`select * from customers where %s;`, details)); err != nil {
		CheckError("Error getting Customer details.", err, false)
		return Customers{}, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&customer.ID, &customer.CardCode, &customer.CardName, &customer.Address, &customer.Phone, &customer.Phone1, &customer.City, &customer.Email, &customer.Synced, &customer.CreatedBy, &customer.CreatedAt, &customer.UpdatedAt, &customer.DeletedAt); err != nil {
			CheckError("Error Scanning Customers.", err, false)
		}
	}

	return
}

// RemoveCustomer deletes a Customers from the database
func RemoveCustomer(id int) (err error) {
	if err = Modify(fmt.Sprintf(`delete from customers where id = %d;`, id)); err != nil {
		CheckError("Error removing Customer(s).", err, false)
		return err
	}

	return nil
}

// NewCustomer creates a new Customer in the database
func NewCustomer(customer map[string]interface{}) (id int, err error) {
	if result, err := MaptoInsert(customer, "customers"); err == nil {
		if id, err = Insert(result); err != nil {
			CheckError("Error inserting the Customer(s).", err, false)
		}
	}

	return
}

// UpdateCustomer updates the Customers  details in the database
func UpdateCustomer(customer map[string]interface{}) (id int, err error) {
	if result, err := MaptoUpdate(customer, "customers", "id"); err == nil {
		if err = Modify(result); err != nil {
			CheckError("Error updating the Customer(s).", err, false)
		}
	}

	return
}
