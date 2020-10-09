package service

import (
	"database/sql"
	"fmt"
)

// TodaysOrders returns an array of Orders placed today
func TodaysOrders() (orders []Orders, err error) {
	var rows *sql.Rows
	if rows, err = Get(`select * from orders where deleted_at is null and cast(created_at as date) = CURDATE() order by created_at desc;`); err != nil {
		CheckError("Error getting Orders.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		order := Orders{}
		if err = rows.Scan(&order.ID, &order.DocEntry, &order.DocNum, &order.Canceled, &order.CardCode, &order.CardName, &order.VatSum, &order.DocTotal, &order.Synced, &order.CreatedBy, &order.CreatedAt, &order.UpdatedAt, &order.DeletedAt); err != nil {
			CheckError("Error Scanning Orders.", err, false)
		} else {
			orders = append(orders, order)
		}
	}

	return
}

// WeekTopSellers returns an array of the top selling Products for the week
func WeekTopSellers() (items []OrderedItems, err error) {
	var rows *sql.Rows
	if rows, err = Get(`select i.id, i.orderid, i.itemcode, i.itemname, sum(i.price) price, 
		sum(i.quantity) quantity, sum(i.discount) discount from ordereditems i 
		inner join orders o on o.id = i.orderid WHERE cast(created_at as date) BETWEEN
		cast(DATE_ADD(CURDATE(), INTERVAL(1 - DAYOFWEEK(CURDATE())) DAY) as date) AND
		cast(DATE_ADD(CURDATE(), INTERVAL(7 - DAYOFWEEK(CURDATE())) DAY) as date) AND
		deleted_At is null GROUP BY i.id, i.orderid, i.itemcode, i.itemname;`); err != nil {
		CheckError("Error getting Top Sellers.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		item := OrderedItems{}
		if err = rows.Scan(&item.ID, &item.OrderID, &item.ItemCode, &item.ItemName, &item.Price, &item.Quantity, &item.Discount); err != nil {
			CheckError("Error Scanning Top Sellers.", err, false)
		} else {
			items = append(items, item)
		}
	}

	return
}

// ReportObject struct
type ReportObject struct {
	Items  []OrderedItems `json:"items"`
	Orders []Orders       `json:"orders"`
}

// Dashboard returns the data for the Dashboard.
func Dashboard() (response ReportObject, err error) {
	response.Orders, err = TodaysOrders()
	response.Items, err = WeekTopSellers()

	return
}

// GetDML returns a report DML Struct.
func GetDML(id int) (report Reports, err error) {
	var rows *sql.Rows
	if rows, err = Get(fmt.Sprintf(`select * from reports where id = %d and deleted_at is null;`, id)); err != nil {
		CheckError("Error getting Reports DML data.", err, false)
		return Reports{}, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&report.ID, &report.Title, &report.Query, &report.CreatedBy, &report.CreatedAt, &report.UpdatedAt, &report.DeletedAt); err != nil {
			CheckError("Error Scanning Report DML.", err, false)
		}
	}

	return
}

// GetReport returns a report object.
func GetReport(id int) (allMaps []map[string]interface{}, err error) {
	var rows *sql.Rows
	var columns []string

	if rows, err = Get(fmt.Sprintf(`SET @sql_query = '';
		select query into @sql_query from reports where id = %d and deleted_at is null;
		PREPARE statement FROM @sql_query;
		EXECUTE statement;
		DEALLOCATE PREPARE statement;`, id)); err != nil {
		CheckError("Error getting Report Data.", err, false)
		return nil, err
	}

	if columns, err = rows.Columns(); err != nil {
		CheckError("Error getting Row columns from Report.", err, false)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		// Dynamic Result rows scanning.
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))

		for i := range values {
			pointers[i] = &values[i]
		}

		err = rows.Scan(pointers...)

		resultMap := make(map[string]interface{})
		for i, val := range values {
			resultMap[columns[i]] = fmt.Sprintf("%s", val)
		}

		// for each database row / record, a map with the column names and row values is added to the allMaps slice
		allMaps = append(allMaps, resultMap)
	}

	return
}

// GetReports returns an array of all created reports.
func GetReports() (reports []Reports, err error) {
	var rows *sql.Rows
	if rows, err = Get(`select * from reports where deleted_at is null order by created_at desc;`); err != nil {
		CheckError("Error getting Report List.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		report := Reports{}
		if err = rows.Scan(&report.ID, &report.Title, &report.Query, &report.CreatedBy, &report.CreatedAt, &report.UpdatedAt, &report.DeletedAt); err != nil {
			CheckError("Error Scanning Report list.", err, false)
		} else {
			reports = append(reports, report)
		}
	}

	return
}

// NewReport creates a new Report in the database
func NewReport(title, query string, creator int) (id int64, err error) {
	insertCommand := fmt.Sprintf(`INSERT INTO reports (title, query, created_by) VALUES ("%s", "%s", "%d");`, title, query, creator)
	if id, err = Insert(insertCommand); err != nil {
		CheckError("Error Registering the new Report.", err, false)
	}

	return
}

// UpdateReport updates the Report in the database
func UpdateReport(title, query string, recordID int) (err error) {
	updateCommand := fmt.Sprintf(`UPDATE reports SET title = "%s", query = "%s" WHERE id = "%d";`, title, query, recordID)
	if err = Modify(updateCommand); err != nil {
		CheckError("Error Updating the requested Report.", err, false)
	}

	return
}

// RemoveReport soft deletes a report from the database
func RemoveReport(id int) (err error) {
	if err = Modify(fmt.Sprintf(`update reports set deleted_at = CURRENT_TIMESTAMP where id = %d;`, id)); err != nil {
		CheckError("Error removing Report.", err, false)
		return err
	}

	return nil
}

// GetTableSchema returns the table schema for the Database
func GetTableSchema() (allMaps []map[string]interface{}, err error) {
	var rows *sql.Rows
	var columns []string

	if rows, err = Get("SELECT table_name, json_arrayagg(json_object('column_name', column_name, 'column_type', column_type, 'column_key', column_key)) as columns from information_schema.columns where table_schema = '" + dbname + "' group by table_name order by table_name;"); err != nil {
		CheckError("Error getting Table Schema.", err, false)
		return nil, err
	}

	if columns, err = rows.Columns(); err != nil {
		CheckError("Error getting Row columns from Schema Query.", err, false)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		// Dynamic Result rows scanning.
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))

		for i := range values {
			pointers[i] = &values[i]
		}

		err = rows.Scan(pointers...)

		resultMap := make(map[string]interface{})
		for i, val := range values {
			resultMap[columns[i]] = fmt.Sprintf("%s", val)
		}

		// for each database row / record, a map with the column names and row values is added to the allMaps slice
		allMaps = append(allMaps, resultMap)
	}

	return
}
