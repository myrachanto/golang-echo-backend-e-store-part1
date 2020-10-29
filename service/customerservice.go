package service

import (
	// "fmt"
	"github.com/myrachanto/astore-v2.0/httperors"
	"github.com/myrachanto/astore-v2.0/model"
	r "github.com/myrachanto/astore-v2.0/repository"
	"github.com/myrachanto/astore-v2.0/support"
)

var (
	Customerservice customerservice = customerservice{}

) 
type customerservice struct {
	
}

func (service customerservice) Create(customer *model.Customer) (*model.Customer, *httperors.HttpError) {
	customer, err1 := r.Customerrepo.Create(customer)
	if err1 != nil {
		return nil, err1
	}
	 return customer, nil

}
func (service customerservice) GetOne(id int) (*model.Customer, *httperors.HttpError) {
	customer, err1 := r.Customerrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return customer, nil
}

func (service customerservice) GetAll(search *support.Search) ([]model.Customer, *httperors.HttpError) {
	
	results, err := r.Customerrepo.GetAll(search)
	if err != nil { 
		return nil, err
	}
	return results, nil 
}

func (service customerservice) Update(id int, customer *model.Customer) (*model.Customer, *httperors.HttpError) {
	customer, err1 := r.Customerrepo.Update(id, customer)
	if err1 != nil {
		return nil, err1
	}
	
	return customer, nil
}
func (service customerservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Customerrepo.Delete(id)
		return success, failure
}
///////deleting a batch////////////////////

//db.Where("age = ?", 20).Delete(&User{})