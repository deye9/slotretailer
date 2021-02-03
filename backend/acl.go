package service

import "database/sql"

// GetRoles returns an array of role Names
func GetRoles() (acls []string, err error) {
	var rows *sql.Rows
	if rows, err = Get(`select distinct rolename from acl order by 1 desc;;`); err != nil {
		CheckError("Error getting Role Names.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var acl string
		if err = rows.Scan(&acl); err != nil {
			CheckError("Error Scanning Customers.", err, false)
		} else {
			acls = append(acls, acl)
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
