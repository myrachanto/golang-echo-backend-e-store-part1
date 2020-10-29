package service

import (
	// "fmt"
	"github.com/myrachanto/astore-v2.0/httperors"
	"github.com/myrachanto/astore-v2.0/model"
	r "github.com/myrachanto/astore-v2.0/repository"
)

var (
	Cartservice cartservice = cartservice{}

) 
type cartservice struct {
	
}

func (service cartservice) Create(cart *model.Cart) (*model.Cart, *httperors.HttpError) {
	cart, err1 := r.Cartrepo.Create(cart)
	if err1 != nil {
		return nil, err1
	}
	 return cart, nil

}
func (service cartservice) GetOne(id int) (*model.Cart, *httperors.HttpError) {
	cart, err1 := r.Cartrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return cart, nil
}

func (service cartservice) GetAll(carts []model.Cart) ([]model.Cart, *httperors.HttpError) {
	carts, err := r.Cartrepo.GetAll(carts)
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func (service cartservice) Update(id int, cart *model.Cart) (*model.Cart, *httperors.HttpError) {
	cart, err1 := r.Cartrepo.Update(id, cart)
	if err1 != nil {
		return nil, err1
	}
	
	return cart, nil
}
func (service cartservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Cartrepo.Delete(id)
		return success, failure
}
///////deleting a batch////////////////////

//db.Where("age = ?", 20).Delete(&User{})