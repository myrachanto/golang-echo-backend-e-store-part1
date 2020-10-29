package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/astore-v2.0/httperors"
	"github.com/myrachanto/astore-v2.0/model"
	"github.com/myrachanto/astore-v2.0/support"
	"github.com/vcraescu/go-paginator" 
	"github.com/vcraescu/go-paginator/adapter"
)
  
var (
	Customerrepo customerrepo = customerrepo{}
)

///curtesy to gorm
type customerrepo struct{}

func (customerRepo customerrepo) Create(customer *model.Customer) (*model.Customer, *httperors.HttpError) {
	if err := customer.Validate(); err != nil {
		return nil, err
	}
	ok := customer.ValidateEmail(customer.Email)
	if !ok {
		return nil, httperors.NewNotFoundError("Your email format is wrong!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&customer)
	IndexRepo.DbClose(GormDB)
	return customer, nil
}
func (customerRepo customerrepo) GetOne(id int) (*model.Customer, *httperors.HttpError) {
	ok := customerRepo.customerUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("customer with that id does not exists!")
	}
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	fmt.Println(customer)
	GormDB.Model(&customer).Where("id = ?", id).First(&customer)
	IndexRepo.DbClose(GormDB)
	
	return &customer, nil
}
func (customerRepo customerrepo) GetAll(search *support.Search) ([]model.Customer,*httperors.HttpError) {
	customers := []model.Customer{} 
	results, err1 := customerRepo.Search(search, customers)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

// func (customerRepo customerrepo) GetAll(search *support.Search) ([]interface{}, *httperors.HttpError) {
// 	customer := model.Customer{}
// 	// customers := []model.Customer{}
// 	// results, err1 := customerRepo.Search(search, customer)
// 	 results, err1 := support.SearchQuery(search, customer)
// 	if err1 != nil {
// 			return nil, err1
// 		}
// 	return results, nil 
// }

func (customerRepo customerrepo) Update(id int, customer *model.Customer) (*model.Customer, *httperors.HttpError) {
	ok := customerRepo.customerUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("customer with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	acustomer := model.Customer{}
	
	GormDB.Model(&acustomer).Where("id = ?", id).First(&acustomer)
	if customer.Name  == "" {
		customer.Name = acustomer.Name
	}
	if customer.Company  == "" {
		customer.Company = acustomer.Company
	}
	if customer.Phone  == "" {
		customer.Phone = acustomer.Phone
	}
	if customer.Email  == "" {
		customer.Email = acustomer.Email
	}
	if customer.Address  == "" {
		customer.Address = acustomer.Address
	}
	GormDB.Model(&customer).Where("id = ?", id).First(&customer).Update(&acustomer)
	
	IndexRepo.DbClose(GormDB)

	return customer, nil
}
func (customerRepo customerrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := customerRepo.customerUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("customer with that id does not exists!")
	}
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	GormDB.Model(&customer).Where("id = ?", id).First(&customer)
	GormDB.Delete(customer)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (customerRepo customerrepo)customerUserExistByid(id int) bool {
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&customer, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (customerRepo customerrepo) Search(Ser *support.Search, customers []model.Customer)([]model.Customer, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	customer := model.Customer{}
	// // invoices := model.Invoice{}
	// fmt.Println(&customer)
	switch(Ser.Search_operator){
	case "all":
		//db.Order("name DESC")
		q := GormDB.Model(&customer).Order(Ser.Column+" "+Ser.Direction).Find(&customers)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&customers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&customers);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&customers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&customers);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&customers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		// order := &Order
		// db.Where("id = ? and status = ?", reqOrder.id, "cart")
		// .Preload("OrderItems").Preload("OrderItems.Item").First(&order)
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&customers);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&customers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&customers);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&customers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&customers);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&customers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&customers);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&customers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&customers);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&customers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Preload("Invoices").Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&customers);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&customers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	// break;
	case "like":
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&customers);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&customers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&customers);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&customers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return customers, nil
}
////////////subject to futher scrutiny/////////////////////////////////
// func (customerRepo customerrepo)paginator(q *gorm.DB, Ser *support.Search, customers []model.Customer) ([]model.Customer, *httperors.HttpError) {
// 	p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
// 	p.SetPage(Ser.Page)
// 	// fmt.Println(Ser.Per_page)
// 	err3 := p.Results(&customers)
// 	if err3 != nil {
// 		return nil, httperors.NewNotFoundError("something went wrong paginating!")
// 	}
// 	return customers, nil
	
// }