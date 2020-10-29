package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/astore-v2.0/httperors"
	"regexp"
)

type Customer struct {
	Name string `gorm:"not null"`
	Company string `gorm:"not null"`
	Phone string `gorm:"not null"`
	Email string `gorm:"not null"`
	Address string `gorm:"not null"`
	Invoices []Invoice `gorm:"foreignKey:CustomerID; not null"`
	CustomerID uint `json:"customerid"`
	gorm.Model
}
func (customer Customer) Validate() *httperors.HttpError{ 
	if customer.Name == "" && len(customer.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if customer.Company == "" && len(customer.Company) < 3 {
		return httperors.NewNotFoundError("Invalid Company")
	}
	if customer.Phone == "" {
		return httperors.NewNotFoundError("Invalid Phone")
	}
	if customer.Email == ""{
		return httperors.NewNotFoundError("Invalid Email")
	}
	if customer.Address == "" && len(customer.Address) > 10 {
		return httperors.NewNotFoundError("Invalid Address")
	}
	return nil
}
func (customer Customer)ValidateEmail(email string) (matchedString bool) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&amp;'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matchedString = re.MatchString(email)
	return matchedString
}