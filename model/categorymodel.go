package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/astore-v2.0/httperors"
)

type Category struct {
	Name string `gorm:"not null" json:"name"` 
	Title string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	// MajorcategoryID uint `json:"majorcategoryid"`
	Majorcategory string ` json:"majorcategory"`
	gorm.Model
}
func (category Category) Validate() *httperors.HttpError{ 
	if category.Name == "" && len(category.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if category.Title == "" && len(category.Title) < 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if category.Description == "" && len(category.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}