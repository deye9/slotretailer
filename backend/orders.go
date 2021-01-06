package service

import (
	"database/sql"
	"fmt"
	"strings"
)

// GetOrder returns a instance belonging to the Order id passed in
func GetOrder(id int) (order Orders, err error) {
	var rows *sql.Rows
	if rows, err = Get(fmt.Sprintf(`select * from orders o 
		inner join ordereditems i on o.id = i.orderid 
		where o.id = %d and o.deleted_at is null;`, id)); err != nil {
		CheckError("Error getting Order data.", err, false)
		return Orders{}, err
	}
	defer rows.Close()

	var items []OrderedItems
	for rows.Next() {
		item := OrderedItems{}

		if err = rows.Scan(&order.ID, &order.DocEntry, &order.DocNum, &order.Canceled, &order.CardCode, &order.CardName, &order.VatSum,
			&order.DocTotal, &order.Synced, &order.CreatedBy, &order.CreatedAt, &order.UpdatedAt, &order.DeletedAt, &order.Comment,
			&order.Returned, &order.DiscountApprovedBy, &item.ID, &item.OrderID, &item.ItemCode, &item.ItemName, &item.Price,
			&item.Quantity, &item.Discount, &item.SerialNumber); err != nil {
			CheckError("Error Scanning Order.", err, false)
		} else {
			items = append(items, item)
		}
	}

	order.Items = items
	return
}

// GetOrders returns an array of Orders
func GetOrders() (orders []Orders, err error) {
	var rows *sql.Rows
	if rows, err = Get(`select * from orders where deleted_at is null order by created_at desc;`); err != nil {
		CheckError("Error getting Orders.", err, false)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		order := Orders{}
		if err = rows.Scan(&order.ID, &order.DocEntry, &order.DocNum, &order.Canceled, &order.CardCode, &order.CardName, &order.VatSum, &order.DocTotal, &order.Synced, &order.CreatedBy, &order.CreatedAt, &order.UpdatedAt, &order.DeletedAt, &order.Comment, &order.Returned, &order.DiscountApprovedBy); err != nil {
			CheckError("Error Scanning Orders.", err, false)
		} else {
			orders = append(orders, order)
		}
	}

	return
}

// RemoveOrder deletes a Order from the database
func RemoveOrder(id int) (err error) {
	if err = Modify(fmt.Sprintf(`update orders set deleted_at = CURRENT_TIMESTAMP where id = %d;`, id)); err != nil {
		CheckError("Error removing Order.", err, false)
		return err
	}

	return nil
}

// NewOrder creates a new Order in the database
func NewOrder(order map[string]interface{}) (err error) {
	master := ""
	detail := ""
	payment := ""
	inventory := ""

	// Get a new reference to the ordered items and remove it from the map.
	payments := order["payments"]
	itemsOrdered := order["items"]
	delete(order, "items")
	delete(order, "payments")

	// Get the master insert query
	if master, err = MaptoInsert(order, "orders"); err != nil {
		CheckError("Error Mapping the Order to SQL.", err, false)
		return err
	}

	master += "SET @last_id := (SELECT LAST_INSERT_ID());"

	// Get the details insert query
	for _, _value := range itemsOrdered.([]interface{}) {
		if detail, err = MaptoInsert(_value.(map[string]interface{}), "ordereditems"); err != nil {
			CheckError("Error Mapping the Ordered Items to SQL.", err, false)
			return err
		}

		// Build out the needed queries
		inventory += fmt.Sprintf(`UPDATE products SET onhand = onhand - %v WHERE itemcode = "%s";`, _value.(map[string]interface{})["quantity"], _value.(map[string]interface{})["itemcode"])
		master += strings.Replace(fmt.Sprintf("%v", detail), `""`, "@last_id", -1)
		master = strings.Replace(master, `"null"`, `""`, -1)
	}

	// Get the payments insert query
	for _, _value := range payments.([]interface{}) {
		if detail, err = MaptoInsert(_value.(map[string]interface{}), "payments"); err != nil {
			CheckError("Error Mapping the Payments to SQL.", err, false)
			return err
		}

		// Build out the needed queries
		payment += strings.Replace(fmt.Sprintf("%v", detail), `""`, "@last_id", -1)
	}

	// Save the Order and Reduce inventory
	if err = Modify(master + payment + inventory); err != nil {
		CheckError("Error creating the Order.", err, false)
		return err
	}

	// Print out invoice & email to the customer.
	return
}

// CreateReturn creates a new returns Document.
func CreateReturn(id int, comment string) (err error) {
	sqlCommand := fmt.Sprintf("UPDATE orders set Comment = '%s', returned = true where ID = %d", comment, id)
	if err = Modify(sqlCommand); err != nil {
		CheckError("Error updating the Order.", err, false)
	}
	return
}

// UpdateOrder updates the Order  details in the database
func UpdateOrder(order map[string]interface{}) (id int, err error) {
	if result, err := MaptoUpdate(order, "orders", "id"); err == nil {
		if err = Modify(result); err != nil {
			CheckError("Error updating the Order.", err, false)
		}
	}

	return
}
