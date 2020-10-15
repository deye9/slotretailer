package service

import (
	"database/sql"
	"fmt"
	"strings"
)

// GetBanks returns an array of Banks
func GetBanks() (banks []Banks, err error) {
	var rows *sql.Rows

	if rows, err = Get(`select * from banks order by name;`); err != nil {
		CheckError("Error getting Banks.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		bank := Banks{}
		if err = rows.Scan(&bank.ID, &bank.Code, &bank.Name); err != nil {
			CheckError("Error Scanning Banks.", err, false)
		} else {
			banks = append(banks, bank)
		}
	}

	return
}

// PaymentOnOrder returns the payment on the order
func PaymentOnOrder(orderID int) (payments []Payments, err error) {
	var rows *sql.Rows

	if rows, err = Get(fmt.Sprintf(`select * from payments where orderid = %d;`, orderID)); err != nil {
		CheckError("Error getting Payments data.", err, false)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		payment := Payments{}
		if err = rows.Scan(&payment.ID, &payment.OrderID, &payment.DocEntry, &payment.DocNum, &payment.Canceled,
			&payment.PaymentType, &payment.PaymentDetails, &payment.Amount); err != nil {
			CheckError("Error Scanning Order.", err, false)
		} else {
			payments = append(payments, payment)
		}
	}

	return
}

// Login grants access to valid users
func Login(email, password string) Users {
	rows, err := Get(fmt.Sprintf(`select id, firstname, lastname, email, isadmin from users where email = "%s" and password = "%s" limit 1;`, email, password))
	if err != nil {
		CheckError("Error logging user in.", err, false)
		return Users{}
	}

	defer rows.Close()
	for rows.Next() {
		var err error
		var user Users

		if err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.IsAdmin); err == nil {
			return user
		}
		CheckError("Error Scanning user.", err, false)
	}

	return Users{}
}

// Search returns the search result
func Search(searchText string) (result []SearchResult, err error) {
	var rows *sql.Rows
	var selectQuery string

	if rows, err = Get(fmt.Sprintf("call searchDB('%s', '%s')", dbname, searchText)); err != nil {
		CheckError("Error getting Search Result", err, false)
		return nil, err
	}
	defer rows.Close()

	for rows.NextResultSet() {
		for rows.Next() {
			data := SearchResult{}

			if err = rows.Scan(&data.Column, &data.Occurrences); err != nil {
				err = rows.Scan(&data.Query)
				CheckError("Error Scanning Order.", err, false)
			} else {
				result = append(result, data)
			}
		}
	}

	for _, val := range result {
		column := strings.Replace(val.Column, "`", "", -1)

		// Split on dot.
		result := strings.Split(column, ".")

		// Build out the queries
		selectQuery += fmt.Sprintf(`SELECT id AS "column", concat('%s = ', %s) AS "occurrences", '%s' as owner FROM %s WHERE %s LIKE "%s";`, result[1], result[1], result[0], result[0], result[1], "%"+searchText+"%")
	}

	// Execute the queries in a batch mode
	if rows, err = Get(selectQuery); err != nil {
		CheckError("Error getting batched Search Result", err, false)
		return nil, err
	}
	defer rows.Close()

	// reset the result array
	result = nil
	for rows.NextResultSet() {
		for rows.Next() {
			data := SearchResult{}

			if err = rows.Scan(&data.Column, &data.Occurrences, &data.Query); err != nil {
				CheckError("Error Scanning Child Search Results.", err, false)
			} else {
				result = append(result, data)
			}
		}
	}

	return
}

// GetAuditLog returns all the audit logs generated btw yesterday and today
func GetAuditLog() (auditLogs []map[string]interface{}, err error) {
	// var rows *sql.Rows
	// select a.id, old_row_data, new_row_data, dml_type, concat(u.firstname, ' ', u.lastname) as 'created_by', dml_timestamp as 'timestamp' from audits a inner join users u on a.dml_created_by = u.id where date(dml_timestamp) between date((subdate(current_date, 1))) and date(CURRENT_DATE);

	return
}
