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
	Transactionrepo transactionrepo = transactionrepo{}
)

///curtesy to gorm
type transactionrepo struct{}

func (transactionRepo transactionrepo) Create(transaction *model.Transaction) (*model.Transaction, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&transaction)
	IndexRepo.DbClose(GormDB)
	return transaction, nil
}
func (transactionRepo transactionrepo) GetOne(id int) (*model.Transaction, *httperors.HttpError) {
	ok := transactionRepo.transactionUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("transaction with that id does not exists!")
	}
	transaction := model.Transaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&transaction).Where("id = ?", id).First(&transaction)
	IndexRepo.DbClose(GormDB)
	
	return &transaction, nil
}

func (transactionRepo transactionrepo) GetAll(transactions []model.Transaction,search *support.Search) ([]model.Transaction, *httperors.HttpError) {
	results, err1 := transactionRepo.Search(search, transactions)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (transactionRepo transactionrepo) Update(id int, transaction *model.Transaction) (*model.Transaction, *httperors.HttpError) {
	ok := transactionRepo.transactionUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("transaction with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	atransaction := model.Transaction{}
	
	GormDB.Model(&atransaction).Where("id = ?", id).First(&atransaction)
	// if transaction.Name  == "" {
	// 	transaction.Name = atransaction.Name
	// }
	// if transaction.Qty  == 0 {
	// 	transaction.Qty = atransaction.Qty
	// }
	// if transaction.Price  == 0 {
	// 	transaction.Price = atransaction.Price
	// }
	
	// if transaction.Discount  == 0 {
	// 	transaction.Discount = atransaction.Discount
	// }
	// if transaction.Tax  == 0 {
	// 	transaction.Tax = atransaction.Tax
	// }
	GormDB.Model(&transaction).Where("id = ?", id).First(&transaction).Update(&atransaction)
	
	IndexRepo.DbClose(GormDB)

	return transaction, nil
}
func (transactionRepo transactionrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := transactionRepo.transactionUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("transaction with that id does not exists!")
	}
	transaction := model.Transaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&transaction).Where("id = ?", id).First(&transaction)
	GormDB.Delete(transaction)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (transactionRepo transactionrepo)transactionUserExistByid(id int) bool {
	transaction := model.Transaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&transaction, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (transactionRepo transactionrepo) Search(Ser *support.Search, transactions []model.Transaction)([]model.Transaction, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	transaction := model.Transaction{}
	switch(Ser.Search_operator){
	case "all":
		q := GormDB.Model(&transaction).Order(Ser.Column+" "+Ser.Direction).Find(&transactions)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&transactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&transactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&transactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&transactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&transactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&transactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&transactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&transactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&transactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	// break;
	case "like":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&transactions);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&transactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&transactions);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&transactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return transactions, nil
}