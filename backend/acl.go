package service

import (
	"database/sql"
	"fmt"
)

// GetRoles returns an array of role Names
func GetRoles() (acls []string, err error) {
	var rows *sql.Rows
	if rows, err = Get("select distinct rolename from acl order by 1 desc;"); err != nil {
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

// GetRoleByName returns an object of the current Role
func GetRoleByName(name string) (acls []ACL, err error) {
	var rows *sql.Rows
	if rows, err = Get(fmt.Sprintf(`select * from acl where rolename = '%s' order by menuname;`, name)); err != nil {
		CheckError("Error getting Role Names.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		acl := ACL{}
		if err = rows.Scan(&acl.ID, &acl.RoleName, &acl.MenuName, &acl.CanCreate, &acl.CanUpdate, &acl.CanDelete, &acl.CanView); err != nil {
			CheckError("Error Scanning Role Details.", err, false)
		} else {
			acls = append(acls, acl)
		}
	}

	return
}

// SaveRole returns nil else an error if there is an error
func SaveRole(acl []interface{}) (err error) {
	sqlCmd := ""
	for _, e := range acl {
		data, _ := e.(map[string]interface{})
		if data["id"].(float64) == 0.0 {
			// insert command
			sqlCmd += fmt.Sprintf("INSERT INTO `acl` (`rolename`, `menuname`, `cancreate`, `canupdate`, `candelete`, `canview`) VALUES ('%s', '%s', %t, %t, %t, %t);", data["rolename"], data["menuname"], data["cancreate"], data["canupdate"], data["candelete"], data["canview"])
		} else {
			// update command
			sqlCmd += fmt.Sprintf("UPDATE `acl` SET `rolename` = '%s', `menuname` = '%s', `cancreate` = %t, `canupdate` = %t, `candelete` = %t, `canview` = %t WHERE `id` = %v;", data["rolename"], data["menuname"], data["cancreate"], data["canupdate"], data["candelete"], data["canview"], data["id"])
		}
	}

	// Set the permissions as requested.
	if err = Modify(sqlCmd); err != nil {
		CheckError("Error setting the permissions.", err, false)
		return err
	}

	return
}
