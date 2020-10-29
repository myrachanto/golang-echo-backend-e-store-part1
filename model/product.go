package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/astore-v2.0/httperors"
)


type Product struct {
	Name string `gorm:"not null" json:"name"`
	Title string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	//Subcategory Subcategory `gorm:"foreignKey:UserID; not null"`
	Category string ` json:"category"`
	Majorcategory string ` json:"majorcategory"`
	Picture string `json:"picture"`
	Price string `json:"price"`
	gorm.Model
}
func (product Product) Validate() *httperors.HttpError{ 
	if product.Name == "" && len(product.Name) > 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if product.Title == "" && len(product.Title) > 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if product.Description == "" && len(product.Description) > 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}