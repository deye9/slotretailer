package service

import (
	"database/sql"
	"fmt"
	"strings"
)

// GetOrder returns a instance belonging to the Order id passed in
func GetOrder(id int) (order Orders, err error) {
	var rows *sql.Rows
	if rows, err = Get(fmt.Sprintf(`select * from orders where id = %d;`, id)); err != nil {
		CheckError("Error getting Order data.", err, false)
		return Orders{}, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&order.ID, &order.DocEntry, &order.DocNum, &order.Canceled, &order.CardCode, &order.CardName, &order.VatSum, &order.DocTotal, &order.Synced, &order.CreatedBy, &order.CreatedAt, &order.UpdatedAt, &order.DeletedAt); err != nil {
			CheckError("Error Scanning Order.", err, false)
		}
	}

	return
}

// GetOrders returns an array of Orders
func GetOrders() (orders []Orders, err error) {
	var rows *sql.Rows
	if rows, err = Get(`select * from orders where deleted_at is null;`); err != nil {
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

// RemoveOrder deletes a Order from the database
func RemoveOrder(id int) (err error) {
	if err = Modify(fmt.Sprintf(`delete from orders where id = %d; delete from ordereditems where orderid = %d`, id, id)); err != nil {
		CheckError("Error removing Order.", err, false)
		return err
	}

	return nil
}

// NewOrder creates a new Order in the database
func NewOrder(order map[string]interface{}) (err error) {
	OrderID := 0
	// Get a new reference to the ordered items and remove it from the map.
	itemsOrdered := order["items"]
	delete(order, "items")

	// Insert the sales order
	if result, err := MaptoInsert(order, "orders"); err == nil {
		if OrderID, err = Insert(result); err != nil {
			CheckError("Error inserting the Orders.", err, false)
			return err
		}
	} else {
		CheckError("Error Mapping the Orders.", err, false)
		return err
	}

	// Insert the sales order items
	for _, _value := range itemsOrdered.([]interface{}) {
		if result, err := MaptoInsert(_value.(map[string]interface{}), "ordereditems"); err == nil {
			if _, err = Insert(strings.ReplaceAll(result, "''", string(OrderID))); err != nil {
				CheckError("Error inserting the Ordered Items.", err, false)
				return err
			}
		} else {
			CheckError("Error Mapping the Ordered Items.", err, false)
			return err
		}
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
