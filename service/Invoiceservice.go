package service

import (
	// "fmt"
	"github.com/myrachanto/astore-v2.0/httperors"
	"github.com/myrachanto/astore-v2.0/model"
	r "github.com/myrachanto/astore-v2.0/repository"
	"github.com/myrachanto/astore-v2.0/support"
)

var (
	Invoiceservice invoiceservice = invoiceservice{}

) 
type invoiceservice struct {
	
}

func (service invoiceservice) Create(invoice *model.Invoice) (*model.Invoice, *httperors.HttpError) {
	invoice, err1 := r.Invoicerepo.Create(invoice)
	if err1 != nil {
		return nil, err1
	}
	 return invoice, nil

}
func (service invoiceservice) View() (string, *httperors.HttpError) {
	code, err1 := r.Invoicerepo.View()
	if err1 != nil {
		return "", err1
	}
	return code, nil
}
func (service invoiceservice) GetOne(id int) (*model.Invoice, *httperors.HttpError) {
	invoice, err1 := r.Invoicerepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return invoice, nil
}

func (service invoiceservice) GetAll(invoices []model.Invoice,search *support.Search) ([]model.Invoice, *httperors.HttpError) {
	invoices, err := r.Invoicerepo.GetAll(invoices,search)
	if err != nil {
		return nil, err
	}
	return invoices, nil
}

func (service invoiceservice) Update(id int, invoice *model.Invoice) (*model.Invoice, *httperors.HttpError) {
	invoice, err1 := r.Invoicerepo.Update(id, invoice)
	if err1 != nil {
		return nil, err1
	}
	
	return invoice, nil
}
func (service invoiceservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Invoicerepo.Delete(id)
		return success, failure
}
