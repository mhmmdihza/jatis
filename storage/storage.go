package storage

import (
	"database/sql"
	"reflect"
	"strings"
)

type Customers struct {
	CustomerID          int    `db:"customer_id" columnType:"int" pk:"y"`
	CompanyName         string `db:"company_name" columnType:"varchar(50)"`
	FirstName           string `db:"first_name" columnType:"varchar(30)"`
	LastName            string `db:"last_name" columnType:"varchar(50)"`
	BillingAddress      string `db:"billing_address" columnType:"varchar(255)"`
	City                string `db:"city" columnType:"varchar(50)"`
	StateOrProvince     string `db:"state_or_province" columnType:"varchar(20)"`
	ZIPCode             string `db:"zip_code" columnType:"varchar(20)"`
	Email               string `db:"email" columnType:"varchar(75)"`
	CompanyWebsite      string `db:"company_website" columnType:"varchar(200)"`
	PhoneNumber         string `db:"phone_number" columnType:"varchar(30)"`
	FaxNumber           string `db:"fax_number" columnType:"varchar(30)"`
	ShipAddress         string `db:"ship_address" columnType:"varchar(255)"`
	ShipCity            string `db:"ship_city" columnType:"varchar(50)"`
	ShipStateOrProvince string `db:"ship_state_or_province" columnType:"varchar(50)"`
	ShipZipCode         string `db:"ship_zip_code" columnType:"varchar(20)"`
	ShipPhoneNumber     string `db:"ship_phone_number" columnType:"varchar(30)"`
}
type Employees struct {
	EmployeeID int    `db:"employee_id" columnType:"int" pk:"y"`
	FirstName  string `db:"first_name" columnType:"varchar(50)"`
	LastName   string `db:"last_name" columnType:"varchar(50)"`
	Title      string `db:"title" columnType:"varchar(50)"`
	WorkPhone  string `db:"work_phone" columnType:"varchar(30)"`
}
type ShippingMethods struct {
	ShippingMethodID int    `db:"shipping_method_id" columnType:"int" pk:"y"`
	ShippingMethod   string `db:"shipping_method" columnType:"varchar(20)"`
}
type Orders struct {
	OrderID             int    `db:"order_id" columnType:"int" pk:"y"`
	CustomerID          int    `db:"customer_id" columnType:"int" fk:"Customers(customer_id)"`
	EmployeeID          int    `db:"employee_id" columnType:"int" fk:"Employees(employee_id)"`
	OrderDate           string `db:"order_date" columnType:"date"`
	PurchaseOrderNumber string `db:"purchase_order_number" columnType:"varchar(30)"`
	ShipDate            string `db:"ship_date" columnType:"date"`
	ShippingMethodID    int    `db:"shipping_method_id" columnType:"int" fk:"ShippingMethods(shipping_method_id)"`
	FreightCharge       int    `db:"freight_charge" columnType:"int"`
	Taxes               int    `db:"taxes" columnType:"int"`
	PaymentReceived     string `db:"payment_received" columnType:"char(1)"`
	Comment             string `db:"comment" columnType:"varchar(150)"`
}
type OrderDetails struct {
	OrderDetailID int `db:"order_detail_id" columnType:"int" pk:"y"`
	OrderID       int `db:"order_id" columnType:"int" fk:"Orders(order_id)"`
	ProductID     int `db:"product_id" columnType:"int" fk:"Products(product_id)"`
	Quantity      int `db:"quantity" columnType:"int"`
	UnitPrice     int `db:"unit_price" columnType:"int"`
	Discount      int `db:"discount" columnType:"int"`
}

type Products struct {
	ProductID   int    `db:"product_id" columnType:"int" pk:"y"`
	ProductName string `db:"product_name" columnType:"varchar(50)"`
	UnitPrice   int    `db:"unit_price" columnType:"int"`
	InStock     string `db:"in_stock" columnType:"char(1)"`
}

type Svc struct {
	db *sql.DB
}

func New(db *sql.DB) *Svc {
	dbDdl(db, Customers{}, Employees{}, ShippingMethods{}, Orders{}, OrderDetails{}, Products{})
	return &Svc{db: db}
}

func dbDdl(db *sql.DB, table ...interface{}) {
	var sb strings.Builder
	var fk strings.Builder
	for _, v := range table {
		var pk strings.Builder
		sb.WriteString(" CREATE TABLE IF NOT EXISTS ")
		t := reflect.TypeOf(v)
		tableName := t.Name()
		sb.WriteString(tableName)
		sb.WriteString(" \n( ")
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			sb.WriteString(field.Tag.Get("db"))
			sb.WriteString(" ")
			sb.WriteString(field.Tag.Get("columnType"))
			if field.Tag.Get("pk") == "y" {
				sb.WriteString(" NOT NULL AUTO_INCREMENT ")
				pk.WriteString(" \n PRIMARY KEY (" + field.Tag.Get("db") + ") \n ")
			}
			if field.Tag.Get("fk") != "" {
				fk.WriteString(" ALTER TABLE " + tableName + " ADD FOREIGN KEY IF NOT EXISTS (" + field.Tag.Get("db") + ") REFERENCES " + field.Tag.Get("fk") + "; \n")
			}
			if i != t.NumField()-1 {
				sb.WriteString(" , ")
			} else {
				if pk.Len() != 0 {
					sb.WriteString(" , ")
				}
				sb.WriteString(pk.String())
			}
			sb.WriteString("\n")

		}
		sb.WriteString(" );\n ")
	}
	sb.WriteString(fk.String())
	arr := strings.Split(sb.String(), ";")
	for _, v := range arr {
		blank := strings.TrimSpace(v) == ""
		if blank {
			continue
		}
		_, err := db.Exec(v)
		if err != nil {
			panic(err)
		}
	}

}
