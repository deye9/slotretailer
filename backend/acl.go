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

// GetUserACL returns an array of the users permission
func GetUserACL(id int) (acls []map[string]interface{}, err error) {
	var rows *sql.Rows
	if rows, err = Get(fmt.Sprintf(`select rolename, menuname, cancreate, canupdate, candelete, canview from acl where rolename = (select a.rolename from users u inner join acl a on u.role = a.id where u.id = %d);`, id)); err != nil {
		CheckError("Error getting User Access Details.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var rolename, menuname string
		var cancreate, canupdate, candelete, canview bool

		if err = rows.Scan(&rolename, &menuname, &cancreate, &canupdate, &candelete, &canview); err != nil {
			CheckError("Error Scanning User Access Details.", err, false)
		} else {
			userACL := make(map[string]interface{})
			userACL["canview"] = canview
			userACL["rolename"] = rolename
			userACL["menuname"] = menuname
			userACL["cancreate"] = cancreate
			userACL["canupdate"] = canupdate
			userACL["candelete"] = candelete
			acls = append(acls, userACL)
		}
	}

	return
}

// GetRoleswithID returns an array of role Names alongside their unique ID's
func GetRoleswithID() (acls []map[string]interface{}, err error) {
	var rows *sql.Rows
	if rows, err = Get("select distinct rolename, (select id from acl where rolename = a.rolename limit 1) id from acl a order by 1 desc;"); err != nil {
		CheckError("Error getting Role Names.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var rolename, id string
		if err = rows.Scan(&rolename, &id); err != nil {
			CheckError("Error Scanning Customers.", err, false)
		} else {
			mapD := make(map[string]interface{})
			mapD["rolename"] = rolename
			mapD["id"] = id
			acls = append(acls, mapD)
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
