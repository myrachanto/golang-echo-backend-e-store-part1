package model

import (
	"github.com/jinzhu/gorm"
)
type Transaction struct {
	Productname string 
	Productid uint
	Invoice   Invoice `gorm:"foreignKey:InvoiceID"`
	Code string
	Title string
	Qty uint16 
	Price float64 
	Tax float64 
	Subtotal float64 
	Discount float64 
	Total float64 
	AmountPaid float64 
	Balance float64
	Status bool 
	gorm.Model
}