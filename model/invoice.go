package model

import (
	"time"
	"github.com/jinzhu/gorm"
)
type Invoice struct {
	Customer_id int64 
	Code string
	Title  string
	Description string
	Dated time.Time
	Due_date *time.Time
	Sub_total float64 
	Discount float64 
	Tax float64 
	Total float64 
	PaidStatus bool 
	AmountPaid float64 
	Balance float64
	status bool
	Cn bool
	Transactions []Transaction `gorm:"foreignKey:TransactionID; not null"`
	TransactionID uint `json:"transactionid"`
	gorm.Model
}
type InvoiceItems struct {
	Description string
	Qty int
	Price float64
}