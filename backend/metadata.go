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

// GetAuditLogs returns all the audit logs generated btw yesterday and today
func GetAuditLogs() (auditLogs []AuditLog, err error) {
	var rows *sql.Rows
	if rows, err = Get(`select  id, old_row_data, new_row_data, dml_type, ifnull((select concat(u.firstname, ' ', u.lastname)  from users u where a.dml_created_by = u.id) , 'System Action') as 'created_by', dml_timestamp as 'timestamp' from audits a where date(dml_timestamp) between date((subdate(current_date, 1))) and date(CURRENT_DATE());`); err != nil {
		CheckError("Error getting Audit Logs.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		log := AuditLog{}
		if err = rows.Scan(&log.ID, &log.OldRowData, &log.NewRowData, &log.DmlType, &log.CreatedBy, &log.LoggedOn); err != nil {
			CheckError("Error Scanning Audit Logs data.", err, false)
		} else {
			auditLogs = append(auditLogs, log)
		}
	}
	return
}

// GetAuditLog returns all the audit logs generated btw yesterday and today
func GetAuditLog(startDate, endDate string) (auditlog []AuditLog, err error) {
	var rows *sql.Rows
	if rows, err = Get(fmt.Sprintf(`select  id, old_row_data, new_row_data, dml_type, ifnull((select concat(u.firstname, ' ', u.lastname)  from users u where a.dml_created_by = u.id) , 'System Action') as 'created_by', dml_timestamp as 'timestamp' from audits a where date(dml_timestamp) between date('%s') and date('%s');`, startDate, endDate)); err != nil {
		CheckError("Error getting Audit Log.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		log := AuditLog{}
		if err = rows.Scan(&log.ID, &log.OldRowData, &log.NewRowData, &log.DmlType, &log.CreatedBy, &log.LoggedOn); err != nil {
			CheckError("Error Scanning Audit Log.", err, false)
		} else {
			auditlog = append(auditlog, log)
		}
	}
	return
}
