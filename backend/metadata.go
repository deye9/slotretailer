package service

import (
	"database/sql"
	"fmt"
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
	if rows, err = Get(fmt.Sprintf("call searchDB('%s', '%s')", "retail", searchText)); err != nil {
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

	return
}
