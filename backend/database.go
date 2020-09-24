package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DbConn returns a pointer to the current DBConnection object
var DbConn *sql.DB
var username string = "root"
var password string = "slot"
var dbname string = "retail"
var hostname string = "127.0.0.1:3306"

// Build out the DSN to the database.
func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

// Setup performs all needed operations
func Setup() {
	var err error
	if DbConn, err = connectDB(); err != nil {
		CheckError("Error Connecting to the Database", err, true)
	}

	createDB()
	runMigrations()

	DbConn.SetMaxOpenConns(3)
	DbConn.SetMaxIdleConns(3)
	DbConn.SetConnMaxLifetime(60 * time.Second)
}

// createDB creates a database in the database server
func createDB() {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := DbConn.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	CheckError("Error when creating DB.", err, true)

	_, err = res.RowsAffected()
	CheckError("Error when fetching rows from DB.", err, true)
}

// connectDB returns a *sql.DB instance or an error if present
func connectDB() (*sql.DB, error) {
	return sql.Open("mysql", dsn())
}

// StructToMap Converts a struct to a map while maintaining the json alias as keys
func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		CheckError("StructToMap Conversion Error.", err, false)
	}

	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}

// RunMigrations executes all outstanding migrations.
func runMigrations() (bool, error) {
	var err error = nil
	var fileContent []byte

	if fileContent, err = ReadFile(BasePath() + "/build/migrations.sql"); err != nil {
		CheckError("Error Reading Migrations File: ", err, false)
		return false, err
	}

	if _, err = DbConn.Exec(string(fileContent)); err != nil {
		CheckError("Error Migrating File: ", err, false)
		return false, err
	}

	return true, nil
}

// structToReplace converts a struct to a sqllite3 replace statement
func structToReplace(value interface{}, tableName string) string {
	cmd := ""
	var ok bool
	var listSlice []interface{}

	if listSlice, ok = value.([]interface{}); !ok {
		fmt.Println("Argument is not a slice")
	}

	for _, _value := range listSlice {
		_Value := _value.(map[string]interface{})
		if cmdkey, err := MaptoReplace(_Value, tableName); err != nil {
			CheckError("Error converting struct to replace statement", err, false)
		} else {
			cmd += cmdkey + "\n"
		}
	}

	return cmd + "\n"
}

// MaptoReplace converts a map to a sql replace statement
func MaptoReplace(mapData map[string]interface{}, tableName string) (string, error) {
	cmd := "REPLACE INTO " + tableName
	keys := ""
	values := ""

	for key, value := range mapData {
		keys += fmt.Sprintf("%s,", key)
		newData := strings.Replace(fmt.Sprintf("%v", value), "'", "''", -1)
		values += strings.Replace(fmt.Sprintf("'%v',", newData), "<nil>", "", -1)
	}

	return cmd + "(" + strings.TrimSuffix(keys, ",") + ") VALUES (" + strings.TrimSuffix(values, ",") + ");", nil
}

// MaptoInsert converts a map to a sql insert statement
func MaptoInsert(mapData map[string]interface{}, tableName string) (string, error) {
	cmd := "INSERT INTO " + tableName
	keys := ""
	values := ""

	for key, value := range mapData {
		keys += fmt.Sprintf("%s,", key)
		values += strings.Replace(fmt.Sprintf("'%v',", value), "<nil>", "", -1)
	}

	return cmd + "(" + strings.TrimSuffix(keys, ",") + ") VALUES (" + strings.TrimSuffix(values, ",") + ")", nil
}

// MaptoUpdate converts a map to a sql update statement
func MaptoUpdate(mapData map[string]interface{}, tableName, tableKey string) (string, error) {
	var tblID string = ""
	cmd := "UPDATE " + tableName + " SET "

	for key, value := range mapData {
		if key != tableKey {
			cmd += fmt.Sprintf("%s = ", key)
			cmd += strings.Replace(fmt.Sprintf("'%v',", value), "<nil>", "", -1)
			continue
		}
		tblID = strings.Replace(fmt.Sprintf("'%v'", value), "<nil>", "", -1)
	}

	return strings.TrimSuffix(cmd, ",") + " WHERE " + tableKey + " = " + tblID + ";", nil
}

// Get retrieves data from the data store.
func Get(selectQuery string) (rows *sql.Rows, err error) {
	return DbConn.Query("USE " + dbname + "; " + selectQuery)
}

// Insert creates record(s) in the data store.
func Insert(insertQuery string) (id int, err error) {
	tx, _ := DbConn.Begin()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err = DbConn.ExecContext(ctx, "USE "+dbname+"; "+insertQuery); err == nil {
		row := tx.QueryRow("select last_insert_rowid()") // SQLite specific
		err = row.Scan(&id)
	}

	tx.Commit()
	return
}

// Modify removes / updates record(s) from the data store.
func Modify(modificationQuery string) (err error) {
	tx, _ := DbConn.Begin()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = DbConn.ExecContext(ctx, "USE "+dbname+"; "+modificationQuery)
	tx.Commit()
	return
}

// // structToReplace converts a struct to a sqllite3 replace statement
// func structToReplace_old(value interface{}, tableName string) string {
// 	var cmdkey string
// 	var cmdvalue string

// 	// Print struct with field names and values.
// 	v := reflect.ValueOf(value)
// 	t := v.Type()
// 	for i := 0; i < t.NumField(); i++ {
// 		name := t.Field(i).Name
// 		cmdkey += fmt.Sprintf("%s, ", name)
// 		cmdvalue += fmt.Sprintf("'%v', ", v.FieldByName(name).Interface())
// 	}

// 	return "REPLACE INTO " + tableName + " (" + strings.TrimSuffix(cmdkey, ", ") + ") VALUES (" + strings.TrimSuffix(cmdvalue, ", ") + ");\n"
// }
