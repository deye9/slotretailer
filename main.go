package main

import (
	service "retailmanager/backend"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

func main() {

	// Global Settings
	service.Setup()
	go service.Sync()

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:     1024,
		Height:    768,
		Title:     "RetailManager",
		JS:        js,
		CSS:       css,
		Colour:    "#131313",
		Resizable: true,
	})

	app.Bind(service.Login)

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
	app.Bind(service.GetProducts)
	app.Bind(service.ProductDetails)

	// Metadata Collection
	app.Bind(service.GetBanks)

	// Admin
	app.Bind(service.RemoveOrder)
	app.Bind(service.UpdateOrder)
	app.Bind(service.RemoveCustomer)

	app.Run()
}
