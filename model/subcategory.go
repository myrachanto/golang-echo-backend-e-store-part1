package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/astore-v2.0/httperors"
)

type Subcategory struct {
	Name string `gorm:"not null"`
	Title string `gorm:"not null"`
	Description string `gorm:"not null"`
	Category []Category
	gorm.Model
}
func (subcategory Subcategory) Validate() *httperors.HttpError{ 
	if subcategory.Name == "" && len(subcategory.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if subcategory.Title == "" && len(subcategory.Title) < 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if subcategory.Description == "" && len(subcategory.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}