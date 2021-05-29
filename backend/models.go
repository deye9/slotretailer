package service

import (
	"database/sql"
)

// Products struct
type Products struct {
	ID           string          `json:"id,omitempty"`
	ItemCode     string          `json:"itemcode,omitempty"`
	ItemName     string          `json:"itemname,omitempty"`
	CodeBars     string          `json:"codebars,omitempty"`
	OnHand       int             `json:"onhand,omitempty"`
	MinLevel     int             `json:"minlevel,omitempty"`
	Warehouse    string          `json:"warehouse,omitempty"`
	SerialNumber string          `json:"serialnumbers,omitempty"`
	Manufacturer string          `json:"manufacturer,omitempty"`
	Price        float32         `json:"price,omitempty"`
	Vat          sql.NullFloat64 `json:"vat,omitempty"`
	ItemID       int             `json:"itemid,omitempty"`
}

// CreditCards struct
type CreditCards struct {
	ID   string `json:"id,omitempty"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// CashAccounts struct
type CashAccounts struct {
	ID   string `json:"id,omitempty"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// Cheques struct
type Cheques struct {
	ID   string `json:"id,omitempty"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// BankTransfer struct
type BankTransfer struct {
	ID   string `json:"id,omitempty"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// PriceList struct
type PriceList struct {
	ID   string `json:"id,omitempty"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// Customers struct
type Customers struct {
	ID        int          `json:"id,omitempty"`
	CardCode  string       `json:"cardcode,omitempty"`
	CardName  string       `json:"cardname,omitempty"`
	Address   string       `json:"address,omitempty"`
	Phone     string       `json:"phone,omitempty"`
	Phone1    string       `json:"phone1,omitempty"`
	City      string       `json:"city,omitempty"`
	Email     string       `json:"email,omitempty"`
	Synced    bool         `json:"synced,omitempty"`
	CreatedBy int          `json:"created_by,omitempty"`
	CreatedAt sql.NullTime `json:"created_at,omitempty"`
	UpdatedAt sql.NullTime `json:"updated_at,omitempty"`
	DeletedAt sql.NullTime `json:"deleted_at,omitempty"`
}

// Orders struct
type Orders struct {
	ID                 int            `json:"id,omitempty"`
	DocEntry           int            `json:"docentry"`
	DocNum             int            `json:"docnum"`
	Canceled           bool           `json:"canceled"`
	CardCode           string         `json:"cardcode,omitempty"`
	CardName           string         `json:"cardname,omitempty"`
	VatSum             float32        `json:"vatsum,omitempty"`
	DocTotal           float32        `json:"doctotal,omitempty"`
	Synced             bool           `json:"synced"`
	Items              []OrderedItems `json:"items"`
	Payments           []Payments     `json:"payments"`
	Comment            sql.NullString `json:"comment,omitempty"`
	Returned           bool           `json:"returned"`
	DiscountApprovedBy int            `json:"discountapprovedby"`
	CreatedBy          int            `json:"created_by,omitempty"`
	CreatedAt          sql.NullTime   `json:"created_at,omitempty"`
	UpdatedAt          sql.NullTime   `json:"updated_at,omitempty"`
	DeletedAt          sql.NullTime   `json:"deleted_at,omitempty"`
}

// OrderedItems struct
type OrderedItems struct {
	ID           int     `json:"id"`
	OrderID      int     `json:"orderid"`
	ItemCode     string  `json:"itemcode"`
	ItemName     string  `json:"itemname"`
	Price        float32 `json:"price"`
	Quantity     int     `json:"quantity"`
	Discount     float32 `json:"discount"`
	SerialNumber string  `json:"serialnumber"`
}

// Payments Struct
type Payments struct {
	ID             int     `json:"id,omitempty"`
	OrderID        int     `json:"orderid,omitempty"`
	DocEntry       int     `json:"docentry,omitempty"`
	DocNum         int     `json:"docnum,omitempty"`
	Canceled       bool    `json:"canceled"`
	PaymentType    string  `json:"paymenttype"`
	PaymentDetails string  `json:"paymentdetails,omitempty"`
	Amount         float32 `json:"amount"`
}

// Users struct
type Users struct {
	ID          int                      `json:"id,omitempty"`
	FirstName   string                   `json:"firstname,omitempty"`
	LastName    string                   `json:"lastname,omitempty"`
	Email       string                   `json:"email,omitempty"`
	Password    string                   `json:"password,omitempty"`
	Role        int                      `json:"role,omitempty"`
	Permissions []map[string]interface{} `json:"acl"`
	CreatedBy   int                      `json:"created_by,omitempty"`
	CreatedAt   sql.NullTime             `json:"created_at,omitempty"`
	UpdatedAt   sql.NullTime             `json:"updated_at,omitempty"`
	DeletedAt   sql.NullTime             `json:"deleted_at,omitempty"`
}

// Store struct
type Store struct {
	ID               int          `json:"id,omitempty"`
	Name             string       `json:"name,omitempty"`
	Address          string       `json:"address,omitempty"`
	Phone            string       `json:"phone,omitempty"`
	City             string       `json:"city,omitempty"`
	Email            string       `json:"email,omitempty"`
	OrdersAPI        string       `json:"orders,omitempty"`
	ProductsAPI      string       `json:"products,omitempty"`
	CustomersAPI     string       `json:"customers,omitempty"`
	CreditCardAPI    string       `json:"creditcard,omitempty"`
	TransfersAPI     string       `json:"transfers,omitempty"`
	SyncInterval     int          `json:"sync_interval,omitempty"`
	SapKey           string       `json:"sapkey,omitempty"`
	LogRotation      string       `json:"logrotation,omitempty"`
	AllowVAT         bool         `json:"vat"`
	WarehousesAPI    string       `json:"warehouses,omitempty"`
	PricelistAPI     string       `json:"pricelist,omitempty"`
	ProductPriceList string       `json:"productpricelist,omitempty"`
	CashAccountAPI   string       `json:"cashaccount,omitempty"`
	StoreCashAccount string       `json:"storecashaccount,omitempty"`
	BankTransferAPI  string       `json:"banktransfer,omitempty"`
	ChequesAPI       string       `json:"cheques,omitempty"`
	CreatedBy        int          `json:"created_by,omitempty"`
	CreatedAt        sql.NullTime `json:"created_at,omitempty"`
	UpdatedAt        sql.NullTime `json:"updated_at,omitempty"`
	DeletedAt        sql.NullTime `json:"deleted_at,omitempty"`
}

// SearchResult struct
type SearchResult struct {
	Query       string `json:"query,omitempty"`
	Column      string `json:"column,omitempty"`
	Occurrences string `json:"occurrences,omitempty"`
}

// Reports struct
type Reports struct {
	ID        int          `json:"id,omitempty"`
	Title     string       `json:"title,omitempty"`
	Query     string       `json:"query,omitempty"`
	CreatedBy string       `json:"created_by,omitempty"`
	CreatedAt sql.NullTime `json:"created_at,omitempty"`
	UpdatedAt sql.NullTime `json:"updated_at,omitempty"`
	DeletedAt sql.NullTime `json:"deleted_at,omitempty"`
}

// AuditLog struct
type AuditLog struct {
	ID         int            `json:"id,omitempty"`
	CreatedBy  string         `json:"created_by,omitempty"`
	DmlType    string         `json:"dml_type,omitempty"`
	OldRowData sql.NullString `json:"old_row_data,omitempty"`
	NewRowData string         `json:"new_row_data,omitempty"`
	LoggedOn   sql.NullTime   `json:"timestamp,omitempty"`
}

// ReportObject struct
type ReportObject struct {
	Items  []OrderedItems `json:"items"`
	Orders []Orders       `json:"orders"`
}

// Stores struct
type Stores struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
}

// Transfers struct
type Transfers struct {
	ID        int               `json:"id,omitempty"`
	FromWhs   string            `json:"fromwhs"`
	ToWhs     string            `json:"towhs"`
	Comment   string            `json:"comment"`
	Canceled  bool              `json:"canceled,omitempty"`
	Synced    bool              `json:"synced"`
	Status    string            `json:"status"`
	Items     []Transfereditems `json:"items"`
	CreatedBy int               `json:"created_by,omitempty"`
	CreatedAt sql.NullTime      `json:"created_at,omitempty"`
	UpdatedAt sql.NullTime      `json:"updated_at,omitempty"`
	DeletedAt sql.NullTime      `json:"deleted_at,omitempty"`
	DocEntry  int               `json:"docentry,omitempty"`
	DocNum    int               `json:"docnum,omitempty"`
	RequestID int               `json:"requestid,omitempty"`
}

// Transfereditems struct
type Transfereditems struct {
	ID           int    `json:"id"`
	TransferID   int    `json:"transferid"`
	ItemCode     string `json:"itemcode"`
	ItemName     string `json:"itemname"`
	Quantity     int    `json:"quantity"`
	OnHand       int    `json:"onhand"`
	SerialNumber string `json:"serialnumber"`
}

// ACL struct
type ACL struct {
	ID        int    `json:"id"`
	RoleName  string `json:"rolename"`
	MenuName  string `json:"menuname"`
	CanCreate bool   `json:"cancreate"`
	CanUpdate bool   `json:"canupdate"`
	CanDelete bool   `json:"candelete"`
	CanView   bool   `json:"canview"`
}

type emailParamsStruct struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// REMEMBER TO MODIFY THE migrations.sql ONCE MODIFIED via an alter statement.
