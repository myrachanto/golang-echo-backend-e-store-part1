package repository

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/myrachanto/astore-v2.0/httperors"
	"github.com/myrachanto/astore-v2.0/model"
)

var (
	Cartrepo cartrepo = cartrepo{}
)
type Totals struct {
	Discount float64
	Subtotal float64
	Total float64
}
///curtesy to gorm
type cartrepo struct{}
//////////////
////////////TODO user id///////////
/////////////////////////////////////////
func (cartRepo cartrepo) Create(cart *model.Cart) (*model.Cart, *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	// code, err := Invoicerepo.GeneCode()
	// if err != nil {
	// 	return nil, err
	// }
	// cart.Code = code
	code := cart.Code
	ok := Invoicerepo.InvoiceExistByCode(code)
	if ok == true {
		return nil, httperors.NewNotFoundError("That invoice is already saved!")
	}
	GormDB.Create(&cart)
	IndexRepo.DbClose(GormDB)
	return cart, nil
}
func (cartRepo cartrepo) GetOne(id int) (*model.Cart, *httperors.HttpError) {
	ok := cartRepo.cartUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("cart with that id does not exists!")
	}
	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&cart).Where("id = ?", id).First(&cart)
	IndexRepo.DbClose(GormDB)
	
	return &cart, nil
}

func (cartRepo cartrepo) GetAll(carts []model.Cart) ([]model.Cart, *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	cart := model.Cart{}
	GormDB.Model(&cart).Find(&carts)
	
	IndexRepo.DbClose(GormDB)
	if len(carts) == 0 {
		return nil, httperors.NewNotFoundError("No results found!")
	}
	return carts, nil
}

func (cartRepo cartrepo) Update(id int, cart *model.Cart) (*model.Cart, *httperors.HttpError) {
	ok := cartRepo.cartUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("cart with that id does not exists!")
	}
	
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	acart := model.Cart{}
	
	GormDB.Model(&acart).Where("id = ?", id).First(&acart)
	if cart.Name  == "" {
		cart.Name = acart.Name
	}
	if cart.Qty  == 0 {
		cart.Qty = acart.Qty
	}
	if cart.Price  == 0 {
		cart.Price = acart.Price
	}
	
	if cart.Discount  == 0 {
		cart.Discount = acart.Discount
	}
	if cart.Tax  == 0 {
		cart.Tax = acart.Tax
	}
	GormDB.Model(&cart).Where("id = ?", id).First(&cart).Update(&acart)
	
	IndexRepo.DbClose(GormDB)

	return cart, nil
}
func (cartRepo cartrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := cartRepo.cartUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("cart with that id does not exists!")
	}
	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&cart).Where("id = ?", id).First(&cart)
	GormDB.Delete(cart)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (cartRepo cartrepo) DeleteAll(code string) (*httperors.HttpSuccess, *httperors.HttpError) {
	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Where("code = ?", code).Delete(&cart)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (cartRepo cartrepo)cartUserExistByid(id int) bool {
	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&cart, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (cartRepo cartrepo)SumTotal(code string) (Total *Totals, err *httperors.HttpError) {
	carts := []model.Cart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil,httperors.NewNotFoundError("db connection failed!")
	}
	GormDB.Where("code = ?", code).Find(&carts)
	IndexRepo.DbClose(GormDB)
	for _, t := range carts {
		Total.Discount += t.Discount
		Total.Subtotal += t.Subtotal
		Total.Total += t.Total
	}
	return Total,nil
}
func (cartRepo cartrepo)CarttoTransaction(code string) (tr []model.Transaction, err *httperors.HttpError) {
	carts := []model.Cart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil,httperors.NewNotFoundError("db connection failed!")
	}
	GormDB.Where("code = ?", code).Find(&carts)
	IndexRepo.DbClose(GormDB)
	for _, c := range carts {
		trans := model.Transaction{Productname:c.Name,Productid:c.ProductID,Qty: c.Qty,Price: c.Price,Tax:c.Tax, Code:code, Subtotal:c.Subtotal, Discount:c.Discount,Total:c.Total}
		tr = append(tr, trans)
	}
	return tr,nil
}
