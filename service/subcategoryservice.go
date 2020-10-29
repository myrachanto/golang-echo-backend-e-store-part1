package service

import (
	// "fmt"
	"github.com/myrachanto/astore-v2.0/httperors"
	"github.com/myrachanto/astore-v2.0/model"
	r "github.com/myrachanto/astore-v2.0/repository"
	"github.com/myrachanto/astore-v2.0/support"
)

var (
	Subcategoryservice subcategoryservice = subcategoryservice{}

) 
type subcategoryservice struct {
	
}

func (service subcategoryservice) Create(subcategory *model.Subcategory) (*model.Subcategory, *httperors.HttpError) {
	if err := subcategory.Validate(); err != nil {
		return nil, err
	}	
	subcategory, err1 := r.Subcategoryrepo.Create(subcategory)
	if err1 != nil {
		return nil, err1
	}
	 return subcategory, nil

}
func (service subcategoryservice) GetOne(id int) (*model.Subcategory, *httperors.HttpError) {
	subcategory, err1 := r.Subcategoryrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return subcategory, nil
}

func (service subcategoryservice) GetAll(subcategorys []model.Subcategory,search *support.Search) ([]model.Subcategory, *httperors.HttpError) {
	subcategorys, err := r.Subcategoryrepo.GetAll(subcategorys,search)
	if err != nil {
		return nil, err
	}
	return subcategorys, nil
}

func (service subcategoryservice) Update(id int, subcategory *model.Subcategory) (*model.Subcategory, *httperors.HttpError) {
	subcategory, err1 := r.Subcategoryrepo.Update(id, subcategory)
	if err1 != nil {
		return nil, err1
	}
	
	return subcategory, nil
}
func (service subcategoryservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Subcategoryrepo.Delete(id)
		return success, failure
}
