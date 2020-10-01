package service

import "database/sql"

// Products struct
type Products struct {
	ID           string  `json:"id,omitempty"`
	ItemCode     string  `json:"itemcode,omitempty"` // 20
	ItemName     string  `json:"itemname,omitempty"` // 100
	CodeBars     string  `json:"codebars,omitempty"` // 254
	OnHand       int     `json:"onhand,omitempty"`
	MinLevel     int     `json:"minlevel,omitempty"`
	Warehouse    string  `json:"warehouse,omitempty"`    // 8
	SerialNumber string  `json:"serialnumber,omitempty"` // 17
	Manufacturer string  `json:"manufacturer,omitempty"` //
	Price        float32 `json:"price,omitempty"`
	Vat          float32 `json:"vat,omitempty"`
}

// Banks struct
type Banks struct {
	ID   string `json:"id,omitempty"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// Customers struct
type Customers struct {
	ID        int            `json:"id,omitempty"`
	CardCode  string         `json:"cardcode,omitempty"` // 15
	CardName  string         `json:"cardname,omitempty"` // 100
	Address   string         `json:"address,omitempty"`  // 100
	Phone     string         `json:"phone,omitempty"`    // 20
	Phone1    string         `json:"phone1,omitempty"`   // 20
	City      string         `json:"city,omitempty"`     // 100
	Email     string         `json:"email,omitempty"`    // 100
	Synced    bool           `json:"synced,omitempty"`
	CreatedBy int            `json:"created_by,omitempty"`
	CreatedAt sql.NullString `json:"created_at,omitempty"`
	UpdatedAt sql.NullString `json:"updated_at,omitempty"`
	DeletedAt sql.NullString `json:"deleted_at,omitempty"`
}

// Orders struct
type Orders struct {
	ID        int            `json:"id,omitempty"`
	DocEntry  int            `json:"docentry"`
	DocNum    int            `json:"docnum"`
	Canceled  bool           `json:"canceled"`
	CardCode  string         `json:"cardcode,omitempty"` // 15
	CardName  string         `json:"cardname,omitempty"` // 100
	VatSum    float32        `json:"vatsum,omitempty"`
	DocTotal  float32        `json:"doctotal,omitempty"`
	Synced    bool           `json:"synced"`
	Items     []OrderedItems `json:"items"`
	Payments  []Payments     `json:"payments"`
	CreatedBy int            `json:"created_by,omitempty"`
	CreatedAt sql.NullString `json:"created_at,omitempty"`
	UpdatedAt sql.NullString `json:"updated_at,omitempty"`
	DeletedAt sql.NullString `json:"deleted_at,omitempty"`
}

// OrderedItems struct
type OrderedItems struct {
	ID       int     `json:"id"`
	OrderID  int     `json:"orderid"`
	ItemCode string  `json:"itemcode"` // 20
	ItemName string  `json:"itemname"` // 100
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
	Discount float32 `json:"discount"`
}

// Payments Struct
type Payments struct {
	ID             int     `json:"id,omitempty"`
	OrderID        int     `json:"orderid"`
	DocEntry       int     `json:"docentry"`
	DocNum         int     `json:"docnum"`
	Canceled       bool    `json:"canceled"`
	PaymentType    string  `json:"paymenttype"`
	PaymentDetails string  `json:"paymentdetails,omitempty"`
	Amount         float32 `json:"amount"`
}

// Users struct
type Users struct {
	ID        int            `json:"id,omitempty"`
	FirstName string         `json:"firstname,omitempty"` // 100
	LastName  string         `json:"lastname,omitempty"`  // 100
	Email     string         `json:"email,omitempty"`     // 100
	Password  string         `json:"password,omitempty"`  // 100
	IsAdmin   bool           `json:"isadmin"`
	CreatedBy int            `json:"created_by,omitempty"`
	CreatedAt sql.NullString `json:"created_at,omitempty"`
	UpdatedAt sql.NullString `json:"updated_at,omitempty"`
	DeletedAt sql.NullString `json:"deleted_at,omitempty"`
}

// Store struct
type Store struct {
	ID           int            `json:"id,omitempty"`
	Name         string         `json:"name,omitempty"`    // 100
	Address      string         `json:"address,omitempty"` // 100
	Phone        string         `json:"phone,omitempty"`   // 20
	City         string         `json:"city,omitempty"`    // 100
	Email        string         `json:"email,omitempty"`   // 100
	OrdersAPI    string         `json:"orders,omitempty"`
	ProductsAPI  string         `json:"products,omitempty"`
	CustomersAPI string         `json:"customers,omitempty"`
	BanksAPI     string         `json:"banks,omitempty"`
	SyncInterval int            `json:"sync_interval,omitempty"`
	SapKey       string         `json:"sapkey,omitempty"`
	CreatedBy    int            `json:"created_by,omitempty"`
	CreatedAt    sql.NullString `json:"created_at,omitempty"`
	UpdatedAt    sql.NullString `json:"updated_at,omitempty"`
	DeletedAt    sql.NullString `json:"deleted_at,omitempty"`
}

// REMEMBER TO MODIFY THE migrations.sql ONCE MODIFIED via an alter statement.
