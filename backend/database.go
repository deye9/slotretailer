package service

import (
	"bufio"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	// Default mysql Driver
	_ "github.com/go-sql-driver/mysql"
)

// DbConn returns a pointer to the current DBConnection object
var DbConn *sql.DB
var dbname string
var username string
var password string
var hostname string

// Build out the DSN to the database.
func dsn(dbName string) string {
	var err error = nil
	var file io.Reader

	if file, err = os.Open(BasePath() + "/build/.env"); err != nil {
		CheckError("Unable to read the .env for the database connection setup.", err, true)
		return ""
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		var text []string = strings.Split(scanner.Text(), "=")

		switch strings.ToLower(text[0]) {
		case "dbname":
			dbname = text[1]

		case "username":
			username = text[1]

		case "password":
			password = text[1]

		case "hostname":
			hostname = text[1]
		}
	}

	if dbName == "sys" {
		return fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true&parseTime=true", username, password, hostname, "sys")
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true&parseTime=true", username, password, hostname, dbname)
}

// Setup performs all needed operations
func Setup() {
	var err error

	if _, err = os.Stat(BasePath() + "/build/.env"); os.IsNotExist(err) {
		CheckError("No .env file exists for the database connection setup.", err, true)
		return
	}

	// Run migrations and create the DB only if there is a valid migration file.
	if _, err = os.Stat(BasePath() + "/build/migrations.sql"); !os.IsNotExist(err) {
		if DbConn, err = connectDB("sys"); err != nil {
			CheckError("Error Connecting to the Database", err, true)
			return
		}
		createDB()
		runMigrations()
	}

	if DbConn, err = connectDB(""); err != nil {
		CheckError("Error Connecting to the Database", err, true)
		return
	}

	DbConn.SetMaxOpenConns(3)
	DbConn.SetMaxIdleConns(3)
	DbConn.SetConnMaxLifetime(60 * time.Second)
}

// createDB creates a database in the database server
func createDB() {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err := DbConn.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	CheckError("Error when creating DB.", err, true)
}

// connectDB returns a *sql.DB instance or an error if present
func connectDB(dbName string) (*sql.DB, error) {
	return sql.Open("mysql", dsn(dbName))
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
		// return false, err
	}

	// Rename the migration file
	renameFile(BasePath()+"/build/migrations.sql", BasePath()+"/build/migrated.sql")
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

// structToInsertUpdate converts a struct to a insertupdate statement
func structToInsertUpdate(value interface{}, tableName string) string {
	cmd := ""
	var ok bool
	var listSlice []interface{}

	if listSlice, ok = value.([]interface{}); !ok {
		fmt.Println("Argument is not a slice")
	}

	// Get the keys
	var keys []string
	for key := range listSlice[0].(map[string]interface{}) {
		keys = append(keys, key)
	}

	// Sort the keys
	sort.Strings(keys)

	// Handle Orders updating
	if strings.ToLower(tableName) == "orders" {
		for _, _value := range listSlice {
			_Value := _value.(map[string]interface{})
			cmd += fmt.Sprintf(`update orders set docentry = %.0f, docnum = %.0f, cardcode = "%s" where id = %.0f;`,
				_Value["docEntry"], _Value["docNum"], _Value["cardCode"], _Value["id"])
		}

		return cmd + "\n"
	}

	for _, _value := range listSlice {
		var update string
		var cmdValues []string
		_Value := _value.(map[string]interface{})
		for _, val := range keys {

			if strings.ToLower(val) == "id" && fmt.Sprintf("%.0f", _Value["id"]) == "0" {
				cmdValues = append(cmdValues, fmt.Sprintf(`%v, `, "null"))
				continue
			}

			if strings.ToLower(val) == "valid_until" {
				_Value[val] = "2038-01-19 03:14:07"
			}

			if strings.ToLower(val) == "synced" {
				update += fmt.Sprintf(`%s = "%v", `, val, true)
				cmdValues = append(cmdValues, fmt.Sprintf(`"%v", `, true))
				continue
			}

			cleanedValue := strings.ReplaceAll(fmt.Sprintf(`%v`, _Value[val]), `"`, `\"`)
			update += fmt.Sprintf(`%s = "%v", `, val, cleanedValue)
			cmdValues = append(cmdValues, fmt.Sprintf(`"%v", `, cleanedValue))
		}

		// joining the string array by ", " separator
		cmd += "INSERT INTO " + tableName + "(" + strings.Join(keys, ", ") + ") VALUES (" + strings.TrimSuffix(strings.Join(cmdValues, ""), ", ") + ") ON DUPLICATE KEY UPDATE " + strings.TrimSuffix(update, ", ") + ";  \n"
	}

	return cmd + "\n"
}

// structToInsertUpdate converts a struct to a insert statement
func structToInsert(value interface{}, tableName string) string {
	cmd := ""
	var ok bool
	var listSlice []interface{}

	if listSlice, ok = value.([]interface{}); !ok {
		fmt.Println("Argument is not a slice")
	}

	// Get the keys
	var keys []string
	for key := range listSlice[0].(map[string]interface{}) {
		if strings.ToLower(key) != "items" && strings.ToLower(key) != "requestid" {
			keys = append(keys, key)
		}
	}

	// Sort the keys
	sort.Strings(keys)

	for _, _value := range listSlice {
		var cmdValues []string
		_Value := _value.(map[string]interface{})

		for _, val := range keys {
			if strings.ToLower(val) == "synced" {
				cmdValues = append(cmdValues, fmt.Sprintf(`"%v", `, true))
				continue
			}

			if strings.ToLower(val) == "docdate" {
				date := fmt.Sprintf(`"%s", `, _Value[val])
				date = strings.Split(strings.ToLower(date), "t")[0]
				cmdValues = append(cmdValues, date+"\", ")
				continue
			}
			cmdValues = append(cmdValues, fmt.Sprintf(`"%v", `, _Value[val]))
		}

		// joining the string array by ", " separator
		cmd += "INSERT INTO " + tableName + "(" + strings.Join(keys, ", ") + ") VALUES (" + strings.TrimSuffix(strings.Join(cmdValues, ""), ", ") + ");  \n"

		// Convert the items to insert commands
		items := _Value["items"]
		for _, value := range items.([]interface{}) {
			itemsData := value.(map[string]interface{})

			cmd += "INSERT INTO transfereditems(transferid, itemcode, itemname, quantity, serialnumber) VALUES (" +
				fmt.Sprintf("%.0f", _Value["id"]) + "," +
				fmt.Sprintf(`"%v"`, itemsData["itemCode"]) + "," +
				fmt.Sprintf(`"%v"`, itemsData["itemName"]) + "," +
				fmt.Sprintf(`"%v"`, itemsData["quantity"]) + "," +
				fmt.Sprintf(`"%v"`, itemsData["serialNumber"]) + ");  \n"
		}
	}

	return cmd + "\n"
}

// structToUpdate converts a struct to an update statement
func structToUpdate(value interface{}, tableName string) string {
	cmd := ""
	var ok bool
	var listSlice []interface{}

	if listSlice, ok = value.([]interface{}); !ok {
		fmt.Println("Argument is not a slice")
	}

	// Get the keys
	var keys []string
	for key := range listSlice[0].(map[string]interface{}) {
		if strings.ToLower(key) != "items" && strings.ToLower(key) != "id" && strings.ToLower(key) != "docdate" && strings.ToLower(key) != "canceled" && strings.ToLower(key) != "towhs" && strings.ToLower(key) != "fromwhs" {
			keys = append(keys, key)
		}
	}

	// Sort the keys
	sort.Strings(keys)

	for _, _value := range listSlice {
		id := ""
		update := ""
		_Value := _value.(map[string]interface{})
		for _, val := range keys {

			if strings.ToLower(val) == "requestid" {
				id = fmt.Sprintf("%v", _Value["requestId"])
			}

			if strings.ToLower(val) == "synced" {
				update += fmt.Sprintf(`%s = "%v", `, val, true)
				continue
			}

			cleanedValue := strings.ReplaceAll(fmt.Sprintf(`%v`, _Value[val]), `"`, `\"`)
			update += fmt.Sprintf(`%s = "%v", `, val, cleanedValue)
		}

		// Convert the items to insert commands
		items := _Value["items"]
		for _, value := range items.([]interface{}) {
			itemsData := value.(map[string]interface{})

			cmd += "UPDATE transfereditems SET " + fmt.Sprintf(`quantity = "%v"`, itemsData["quantity"]) + ", " +
				fmt.Sprintf(`serialnumber = "%v"`, itemsData["serialNumber"]) + " WHERE itemcode = " + fmt.Sprintf(`"%v"`, itemsData["itemCode"]) +
				" AND transferid =  " + id + "; \n"
		}

		// joining the string array by ", " separator
		cmd += "UPDATE " + tableName + " SET " + strings.TrimSuffix(update, ", ") + " WHERE ID = " + id + "; \n"
	}

	return cmd + "\n"
}

// MaptoInsert converts a map to a sql insert statement
func MaptoInsert(mapData map[string]interface{}, tableName string) (string, error) {
	cmd := "INSERT INTO " + tableName
	keys := ""
	values := ""

	for key, value := range mapData {
		keys += fmt.Sprintf("%s,", key)
		values += strings.Replace(fmt.Sprintf(`"%v",`, value), "<nil>", "", -1)
	}

	return cmd + "(" + strings.TrimSuffix(keys, ",") + ") VALUES (" + strings.TrimSuffix(values, ",") + "); ", nil
}

// MaptoUpdate converts a map to a sql update statement
func MaptoUpdate(mapData map[string]interface{}, tableName, tableKey string) (string, error) {
	var tblID string = ""
	cmd := "UPDATE " + tableName + " SET updated_at = CURRENT_TIMESTAMP, "

	for key, value := range mapData {
		if key != tableKey {
			cmd += fmt.Sprintf("%s = ", key)
			cmd += strings.Replace(fmt.Sprintf(`"%v",`, value), "<nil>", "", -1)
			continue
		}
		tblID = strings.Replace(fmt.Sprintf(`"%v"`, value), "<nil>", "", -1)
	}

	return strings.TrimSuffix(cmd, ",") + " WHERE " + tableKey + " = " + tblID + ";", nil
}

// Get retrieves data from the data store.
func Get(selectQuery string) (rows *sql.Rows, err error) {
	// // DbConn.QueryRow()
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// if rows, err = DbConn.QueryContext(ctx, selectQuery); err != nil {
	// 	return nil, err
	// }

	// return
	return DbConn.Query(selectQuery)
}

// Insert creates record(s) in the data store.
func Insert(insertQuery string) (id int64, err error) {
	insertQuery = strings.Replace(insertQuery, `"true"`, "true", -1)
	insertQuery = strings.Replace(insertQuery, `"false"`, "false", -1)

	var row driver.Result
	tx, _ := DbConn.Begin()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if row, err = DbConn.ExecContext(ctx, insertQuery); err == nil {
		id, err = row.LastInsertId()
	}

	tx.Commit()
	return
}

// Modify removes / updates record(s) from the data store.
func Modify(modificationQuery string) (err error) {
	modificationQuery = strings.Replace(modificationQuery, `"true"`, "true", -1)
	modificationQuery = strings.Replace(modificationQuery, `"false"`, "false", -1)

	// tx, _ := DbConn.Begin()
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// _, err = DbConn.ExecContext(ctx, modificationQuery)
	// tx.Commit()
	_, err = DbConn.Exec(modificationQuery)
	return
}
