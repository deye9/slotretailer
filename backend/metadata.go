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
	// sqlQuery := fmt.Sprintf("SET @table_schema = 'retail';
	// 			SET @condition = 'LIKE \'%%s%\'';
	// 			SET @column_types_regexp = '^((var)?char|(var)?binary|blob|text|enum|set)\\(';

	// 			-- Reset @sql_query in case it was used previously
	// 			SET @sql_query = '';

	// 			-- Build query for each table and merge with previous queries with UNION
	// 			SELECT
	// 				-- Use `DISTINCT IF(QUERYBUILDING, NULL, NULL)`
	// 				-- to only select a single null value
	// 				-- instead of selecting the query over and over again as it's built
	// 				DISTINCT IF(@sql_query := CONCAT(
	// 					IF(LENGTH(@sql_query), CONCAT(@sql_query, ' UNION '), ''),
	// 					'SELECT ',
	// 						QUOTE(CONCAT('`', `table_name`, '`.`', `column_name`, '`')), ' AS `column`, ',
	// 						'COUNT(*) AS `occurrences` ',
	// 					'FROM `', `table_schema`, '`.`', `table_name`, '` ',
	// 					'WHERE `', `column_name`, '` ', @condition
	// 				), NULL, NULL) `query`
	// 			FROM (
	// 				SELECT
	// 					table_schema,
	// 					table_name,
	// 					column_name
	// 				FROM information_schema.columns
	// 				WHERE table_schema = @table_schema
	// 				AND column_type REGEXP @column_types_regexp
	// 			) results;

	// 			-- Only return results with at least one occurrence
	// 			SET @sql_query = CONCAT('SELECT * FROM (', @sql_query, ') results WHERE occurrences > 0');

	// 			-- Run built query
	// 			PREPARE statement FROM @sql_query;
	// 			EXECUTE statement;
	// 			DEALLOCATE PREPARE statement;", searchText)

	// fmt.Printls(sqlQuery)
	return
}
