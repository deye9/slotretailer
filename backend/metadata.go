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
