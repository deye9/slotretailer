package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// GetPaymentDetails returns the data for all accepted Payments
func GetPaymentDetails() (details []map[string]interface{}, err error) {
	var rows *sql.Rows
	m := make(map[string]interface{})

	if rows, err = Get(`select * from cheques order by name; 
						select * from creditcards order by name;
						select * from banktransfer order by name;`); err != nil {
		CheckError("Error getting Payment Details.", err, false)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		if data, err := ConvertToMAP(rows); err == nil {
			m["cheque"] = data
		}
	}

	for rows.NextResultSet() {
		for rows.Next() {
			if data, err := ConvertToMAP(rows); err == nil {
				if v := m["pos"]; v == nil {
					m["pos"] = data
				} else {
					m["banktransfer"] = data
				}
			}
		}
	}

	details = append(details, m)
	return
}

// ConvertToMAP converts the database response to JSON.
func ConvertToMAP(rows *sql.Rows) (allMaps []map[string]interface{}, err error) {
	var columns []string

	if columns, err = rows.Columns(); err != nil {
		return nil, errors.New("An Error executing at method {ConvertToMAP} while getting row columns for passed in rows")
	}

	// Ensure we only get the column information one time!
	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))
	for i := range values {
		pointers[i] = &values[i]
	}

	for rows.Next() {
		// Dynamic Result rows scanning.
		err = rows.Scan(pointers...)

		resultMap := make(map[string]interface{})
		for i, val := range values {
			if strings.ToLower(columns[i]) != "password" && strings.ToLower(columns[i]) != "created_at" && strings.ToLower(columns[i]) != "updated_at" && strings.ToLower(columns[i]) != "deleted_at" {
				resultMap[columns[i]] = fmt.Sprintf("%s", val)
			}
		}

		// for each database row / record, a map with the column names and row values is added to the allMaps slice
		allMaps = append(allMaps, resultMap)
	}
	return
}

// GetCreditcards returns an array of Banks
func GetCreditcards() (banks []CreditCards, err error) {
	var rows *sql.Rows

	if rows, err = Get(`select * from creditcards order by name;`); err != nil {
		CheckError("Error getting Banks.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		bank := CreditCards{}
		if err = rows.Scan(&bank.ID, &bank.Code, &bank.Name); err != nil {
			CheckError("Error Scanning Banks.", err, false)
		} else {
			banks = append(banks, bank)
		}
	}

	return
}

// GetPriceList returns an array of PriceList
func GetPriceList() (pricelists []PriceList, err error) {
	var rows *sql.Rows

	if rows, err = Get(`select * from pricelist order by name;`); err != nil {
		CheckError("Error getting Price List.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		pricelist := PriceList{}
		if err = rows.Scan(&pricelist.ID, &pricelist.Code, &pricelist.Name); err != nil {
			CheckError("Error Scanning Price List.", err, false)
		} else {
			pricelists = append(pricelists, pricelist)
		}
	}

	return
}

// GetCashAccounts returns an array of Cash Accounts
func GetCashAccounts() (cashaccounts []CashAccounts, err error) {
	var rows *sql.Rows

	if rows, err = Get(`select * from CashAccounts order by name;`); err != nil {
		CheckError("Error getting Cash Accounts.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		cashaccount := CashAccounts{}
		if err = rows.Scan(&cashaccount.ID, &cashaccount.Code, &cashaccount.Name); err != nil {
			CheckError("Error Scanning Cash Accounts.", err, false)
		} else {
			cashaccounts = append(cashaccounts, cashaccount)
		}
	}

	return
}

// PaymentOnOrder returns the payment on the order
func PaymentOnOrder(orderID int) (payments []Payments, err error) {
	var rows *sql.Rows

	if rows, err = Get(fmt.Sprintf(`select id, orderid, docentry, docnum, canceled, paymenttype,
		case lower(paymenttype) 
			when 'cash' then 'Cash' 
			when 'pos' then (select name from creditcards where code = paymentdetails) 
			when 'cheque' then  (select name from cheques where code = paymentdetails) 
			when 'bank transfer' then  (select name from banktransfer where code = paymentdetails) end paymentdetails, 
			amount from payments where orderid = %d;`, orderID)); err != nil {
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
