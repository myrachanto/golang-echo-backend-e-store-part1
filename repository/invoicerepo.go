package repository

import (
	"fmt"
	"strings"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/myrachanto/astore-v2.0/httperors"
	"github.com/myrachanto/astore-v2.0/model"
	"github.com/myrachanto/astore-v2.0/support"
	"github.com/vcraescu/go-paginator" 
	"github.com/vcraescu/go-paginator/adapter"
)

var (
	Invoicerepo invoicerepo = invoicerepo{}
)

///curtesy to gorm
type invoicerepo struct{}

func (invoiceRepo invoicerepo) Create(invoice *model.Invoice) (*model.Invoice, *httperors.HttpError) {
	code := invoice.Code
	t,r := Cartrepo.SumTotal(code);if r != nil {
		return nil, r
	}
	invoice.Discount = t.Discount
	invoice.Sub_total = t.Subtotal
	invoice.Total = t.Total
	tr,e := Cartrepo.CarttoTransaction(code);if e != nil {
		return nil, e
	}
	invoice.Transactions = tr 
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&invoice)
	Cartrepo.DeleteAll(code)
	IndexRepo.DbClose(GormDB)
	return invoice, nil
}
func (invoiceRepo invoicerepo) View() (string, *httperors.HttpError) {
	code,err := invoiceRepo.GeneCode()
	return code, err
}
func (invoiceRepo invoicerepo) GetOne(id int) (*model.Invoice, *httperors.HttpError) {
	ok := invoiceRepo.invoiceUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("invoice with that id does not exists!")
	}
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&invoice).Where("id = ?", id).First(&invoice)
	IndexRepo.DbClose(GormDB)
	
	return &invoice, nil
}

func (invoiceRepo invoicerepo) GetAll(invoices []model.Invoice,search *support.Search) ([]model.Invoice, *httperors.HttpError) {
	
	results, err1 := invoiceRepo.Search(search, invoices)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (invoiceRepo invoicerepo) Update(id int, invoice *model.Invoice) (*model.Invoice, *httperors.HttpError) {
	ok := invoiceRepo.invoiceUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("invoice with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	ainvoice := model.Invoice{}
	
	GormDB.Model(&ainvoice).Where("id = ?", id).First(&ainvoice)
	// if invoice.invoice  == "" {
	// 	invoice.invoice = ainvoice.invoice
	// }
	// if invoice.Description  == "" {
	// 	invoice.Description = ainvoice.Description
	// }
	// if invoice.Subtotal  == 0 {
	// 	invoice.Subtotal = ainvoice.Subtotal
	// }
	// if invoice.Discount  == 0 {
	// 	invoice.Discount = ainvoice.Discount
	// }	
	// if invoice.AmountPaid  == 0 {
	// 	invoice.AmountPaid = ainvoice.AmountPaid
	// }
	GormDB.Model(&invoice).Where("id = ?", id).First(&invoice).Update(&ainvoice)
	
	IndexRepo.DbClose(GormDB)

	return invoice, nil
}
func (invoiceRepo invoicerepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := invoiceRepo.invoiceUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("invoice with that id does not exists!")
	}
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&invoice).Where("id = ?", id).First(&invoice)
	GormDB.Delete(invoice)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (invoiceRepo invoicerepo)invoiceUserExistByid(id int) bool {
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&invoice, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (invoiceRepo invoicerepo)InvoiceExistByCode(code string) bool {
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&invoice, "code =?", code).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB) 
	return true
	
}
func (invoiceRepo invoicerepo)GeneCode() (string, *httperors.HttpError) {
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	if GormDB.Last(&invoice).RecordNotFound(){
		c1 := 1
		code := "CustInvNo_"+ string(c1)
		return code, nil
	 }
	GormDB.Last(&invoice)
	c1 := invoice.ID + 1
	code := "CustInvNo_"+ string(c1)
	IndexRepo.DbClose(GormDB)
	return code, nil
	
}
func (invoiceRepo invoicerepo) Search(Ser *support.Search, invoices []model.Invoice)([]model.Invoice, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	invoice := model.Invoice{}
	switch(Ser.Search_operator){
	case "all":
		q := GormDB.Model(&invoice).Order(Ser.Column+" "+Ser.Direction).Find(&invoices)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&invoices); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&invoices); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&invoices); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&invoices); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&invoices); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&invoices); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&invoices); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&invoices); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&invoices); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	// break;
	case "like":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&invoices); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&invoices); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return invoices, nil
}