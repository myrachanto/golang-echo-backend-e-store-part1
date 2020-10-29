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
	Productrepo productrepo = productrepo{}
)

///curtesy to gorm
type productrepo struct{}

func (productRepo productrepo) Create(product *model.Product) (*model.Product, *httperors.HttpError) {
	if err := product.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	fmt.Println(product)
	GormDB.Create(&product)
	IndexRepo.DbClose(GormDB)
	return product, nil
}
func (productRepo productrepo) GetOne(id int) (*model.Product, *httperors.HttpError) {
	ok := productRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("product with that id does not exists!")
	}
	product := model.Product{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&product).Where("id = ?", id).First(&product)
	IndexRepo.DbClose(GormDB)
	
	return &product, nil
}

func (productRepo productrepo) View() ([]model.Category, *httperors.HttpError) {
	mc, e := Categoryrepo.All()
	if e != nil{
		return nil, e
	}
	return mc, nil
}
func (productRepo productrepo) GetProducts(products []model.Product,search *support.Productsearch) ([]model.Product, *httperors.HttpError) {
	results, err1 := productRepo.SearchFront(search, products)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}
func (productRepo productrepo) GetAll(products []model.Product,search *support.Search) ([]model.Product, *httperors.HttpError) {
	results, err1 := productRepo.Search(search, products)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (productRepo productrepo) Update(id int, product *model.Product) (*model.Product, *httperors.HttpError) {
	ok := productRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("product with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	aproduct := model.Product{}
	
	GormDB.Model(&aproduct).Where("id = ?", id).First(&aproduct)
	if product.Name  == "" {
		product.Name = aproduct.Name
	}
	if product.Title  == "" {
		product.Title = aproduct.Title
	}
	if product.Description  == "" {
		product.Description = aproduct.Description
	}
	if product.Category  == "" {
		product.Category = aproduct.Category
	}
	if product.Majorcategory  == "" {
		product.Majorcategory = aproduct.Majorcategory
	}
	GormDB.Model(&aproduct).Where("id = ?", id).Update(&product)
	
	IndexRepo.DbClose(GormDB)

	return product, nil
}
func (productRepo productrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := productRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	product := model.Product{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&product).Where("id = ?", id).First(&product)
	GormDB.Delete(product)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (productRepo productrepo)ProductUserExistByid(id int) bool {
	product := model.Product{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&product, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (productRepo productrepo) Search(Ser *support.Search, products []model.Product)([]model.Product, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	product := model.Product{}
	switch(Ser.Search_operator){
	case "all":
		q := GormDB.Model(&product).Order(Ser.Column+" "+Ser.Direction).Find(&products)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	// break;
	case "like":
		// fmt.Println(Ser.Search_query_1)
		if Ser.Search_query_1 == "all" {
				//db.Order("name DESC")
		q := GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&products)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		fmt.Println(p.Results(&products))
				if err3 := p.Results(&products); err3 != nil {
					return nil, httperors.NewNotFoundError("something went wrong paginating!")
				}

		}else {

			q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&products);
			p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
			p.SetPage(Ser.Page)
			fmt.Println(p.Results(&products,))
			if err3 := p.Results(&products); err3 != nil {
				return nil, httperors.NewNotFoundError("something went wrong paginating!")
			}
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return products, nil
}

func (productRepo productrepo) SearchFront(Ser *support.Productsearch, products []model.Product)([]model.Product, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	product := model.Product{}
	switch(Ser.Search_operator){
	case "all":
		q := GormDB.Model(&product).Order(Ser.Column+" "+Ser.Direction).Find(&products)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		///product = %as% AND price (>/=/<=/>=/between)
		//db.Where("name = ? AND age >= ? ", "myrachanto", "28").Find(&users)
		//db.Where("name LIKE ?", "%a%").Find(&users)
		q := GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		fmt.Println(Ser)
		q := GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&products); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return products, nil
}