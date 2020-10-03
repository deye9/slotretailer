package service

import "database/sql"

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
