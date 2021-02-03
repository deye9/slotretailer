package main

import (
	service "slotretailer/backend"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

func main() {

	// Global Settings
	service.Setup()

	// Get the store details
	service.GetStore()
	go service.Sync()

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:     1280,
		Height:    720,
		Resizable: true,
		Title:     "SLOT Systems Retail Manager.",
		JS:        js,
		CSS:       css,
		Colour:    "#131313",
	})

	// Access Control List
	app.Bind(service.GetRoles)
	app.Bind(service.GetRoleByName)
	app.Bind(service.SaveRole)
	app.Bind(service.UpdateRole)
	app.Bind(service.DeleteRole)

	// Sync Details
	app.Bind(service.GetLog)
	app.Bind(service.GetLogs)

	// Dashboard
	app.Bind(service.Dashboard)

	// Application Users
	app.Bind(service.NewUser)
	app.Bind(service.GetUser)
	app.Bind(service.GetUsers)
	app.Bind(service.UpdateUser)
	app.Bind(service.RemoveUser)

	// Store Details
	app.Bind(service.GetStore)
	app.Bind(service.SaveStore)
	app.Bind(service.GetStores)

	// Customers Details
	app.Bind(service.NewCustomer)
	app.Bind(service.GetCustomer)
	app.Bind(service.GetCustomers)
	app.Bind(service.UpdateCustomer)
	app.Bind(service.GetCustomerbyPhone)

	// Orders Details
	app.Bind(service.GetOrder)
	app.Bind(service.NewOrder)
	app.Bind(service.GetOrders)

	// Product Details
	app.Bind(service.GetProduct)
	app.Bind(service.GetProducts)
	app.Bind(service.CreateReturn)
	app.Bind(service.ProductDetails)
	app.Bind(service.GetStoreProducts)

	// Inventory Transfers
	app.Bind(service.GetSerials)
	app.Bind(service.GetTransfer)
	app.Bind(service.NewTransfer)
	app.Bind(service.GetTransfers)
	app.Bind(service.UpdateTransfer)
	app.Bind(service.RejectTransfer)
	app.Bind(service.RemoveTransfer)
	app.Bind(service.FinalizeTransfer)

	// Metadata Collection / Application Utilities
	app.Bind(service.Login)
	app.Bind(service.Search)
	app.Bind(service.GetPaymentDetails)
	app.Bind(service.GetPriceList)
	app.Bind(service.GetAuditLog)
	app.Bind(service.GetAuditLogs)
	app.Bind(service.PaymentOnOrder)
	app.Bind(service.GetCashAccounts)

	// Admin
	app.Bind(service.RemoveOrder)
	app.Bind(service.UpdateOrder)
	app.Bind(service.RemoveCustomer)

	// Reports
	app.Bind(service.GetDML)
	app.Bind(service.NewReport)
	app.Bind(service.GetReport)
	app.Bind(service.GetReports)
	app.Bind(service.UpdateReport)
	app.Bind(service.RemoveReport)
	app.Bind(service.GetTableSchema)

	app.Run()
}
