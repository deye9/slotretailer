package service

import (
	"database/sql"
	"fmt"
)

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

// GetUser returns a instance belonging to the user id passed in
func GetUser(id int) (user Users, err error) {
	var rows *sql.Rows
	if rows, err = Get(fmt.Sprintf(`select * from users where id = %d;`, id)); err != nil {
		CheckError("Error getting user data.", err, false)
		return Users{}, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.IsAdmin, &user.CreatedBy, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
			CheckError("Error Scanning user.", err, false)
		}
	}

	return
}

// GetUsers returns an array of users
func GetUsers() (users []Users, err error) {
	var rows *sql.Rows
	if rows, err = Get(`select * from users where deleted_at is null and id > 1;`); err != nil {
		CheckError("Error getting users.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		user := Users{}
		if err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.IsAdmin, &user.CreatedBy, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
			CheckError("Error Scanning users.", err, false)
		} else {
			users = append(users, user)
		}
	}

	return
}

// RemoveUser deletes a user from the database
func RemoveUser(id int) (err error) {
	if err = Modify(fmt.Sprintf(`delete from users where id = %d;`, id)); err != nil {
		CheckError("Error removing user.", err, false)
		return err
	}

	return nil
}

// NewUser creates a new user in the database
func NewUser(user map[string]interface{}) (id int64, err error) {
	if result, err := MaptoInsert(user, "users"); err == nil {
		if id, err = Insert(result); err != nil {
			CheckError("Error inserting the user.", err, false)
		}
	}

	return
}

// UpdateUser updates the user  details in the database
func UpdateUser(user map[string]interface{}) (id int, err error) {
	if result, err := MaptoUpdate(user, "users", "id"); err == nil {
		if err = Modify(result); err != nil {
			CheckError("Error updating the user.", err, false)
		}
	}

	return
}
