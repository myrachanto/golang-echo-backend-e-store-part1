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
	Categoryrepo categoryrepo = categoryrepo{}
)


///curtesy to gorm 
type categoryrepo struct{}

func (categoryRepo categoryrepo) Create(category *model.Category) (*model.Category, *httperors.HttpError) {
	if err := category.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&category)
	IndexRepo.DbClose(GormDB)
	return category, nil
}
func (categoryRepo categoryrepo) View() ([]model.Majorcategory, *httperors.HttpError) {
	mc, e := Majorcategoryrepo.All()
	if e != nil{
		return nil, e
	}
	return mc, nil
}
func (categoryRepo categoryrepo) GetOne(id int) (*model.Category, *httperors.HttpError) {
	ok := categoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("category with that id does not exists!")
	}
	category := model.Category{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&category).Where("id = ?", id).First(&category)
	IndexRepo.DbClose(GormDB)
	
	return &category, nil
}

func (categoryRepo categoryrepo) GetAll(categorys []model.Category,search *support.Search) ([]model.Category, *httperors.HttpError) {
	results, err1 := categoryRepo.Search(search, categorys)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}
func (categoryRepo categoryrepo) All() (t []model.Category, r *httperors.HttpError) {

	category := model.Category{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	q := GormDB.Model(&category).Order("name").Find(&t)
	p := paginator.New(adapter.NewGORMAdapter(q), 40)
	p.SetPage(1)

	
	if err3 := p.Results(&t); err3 != nil {
		return nil, httperors.NewNotFoundError("something went wrong paginating!")
	}
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (categoryRepo categoryrepo) Update(id int, category *model.Category) (*model.Category, *httperors.HttpError) {
	ok := categoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("category with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	Category := model.Category{}
	acategory := model.Category{}
	
	GormDB.Model(&Category).Where("id = ?", id).First(&acategory)
	if category.Name  == "" {
		category.Name = acategory.Name
	}
	if category.Title  == "" {
		category.Title = acategory.Title
	}
	if category.Description  == "" {
		category.Description = acategory.Description
	}
	if category.Majorcategory  == "" {
		category.Majorcategory = acategory.Majorcategory
	}
	GormDB.Model(&Category).Where("id = ?", id).Update(&category)
	
	IndexRepo.DbClose(GormDB)

	return category, nil
}
func (categoryRepo categoryrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := categoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	category := model.Category{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&category).Where("id = ?", id).First(&category)
	GormDB.Delete(category)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (categoryRepo categoryrepo)ProductUserExistByid(id int) bool {
	category := model.Category{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&category, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (categoryRepo categoryrepo) Search(Ser *support.Search, categorys []model.Category)([]model.Category, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	category := model.Category{}
	switch(Ser.Search_operator){
	case "all":
		q := GormDB.Model(&category).Order(Ser.Column+" "+Ser.Direction).Find(&categorys)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&categorys);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&categorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&categorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&categorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&categorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&categorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&categorys);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&categorys);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	// break;
case "like":
	// fmt.Println(Ser.Search_query_1)
	if Ser.Search_query_1 == "all" {
			//db.Order("name DESC")
	q := GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&categorys)
	///////////////////////////////////////////////////////////////////////////////////////////////////////
	///////////////find some other paginator more effective one///////////////////////////////////////////
	p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
	p.SetPage(Ser.Page)
	
	fmt.Println(p.Results(&categorys))
			if err3 := p.Results(&categorys); err3 != nil {
				return nil, httperors.NewNotFoundError("something went wrong paginating!")
			}

	}else {

		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&categorys);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		fmt.Println(p.Results(&categorys,))
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	}
break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&categorys);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&categorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return categorys, nil
}