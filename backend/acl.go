package service

import "database/sql"

// GetRoles returns an array of role Names
func GetRoles() (customer Customers, err error) {
	var rows *sql.Rows
	if rows, err = Get(`select * from customers where deleted_at is null order by created_at desc;`); err != nil {
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

// GetRoleByID returns an object of the current Role
func GetRoleByID(id int) {
	return
}

// SaveRole returns
func SaveRole() {
	return
}

// UpdateRole returns
func UpdateRole() {
	return
}

// DeleteRole returns
func DeleteRole() {
	return
}
