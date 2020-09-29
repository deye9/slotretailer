package service

import "database/sql"

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
