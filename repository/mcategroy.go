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
	Majorcategoryrepo majorcategoryrepo = majorcategoryrepo{}
)

///curtesy to gorm
type majorcategoryrepo struct{}

func (majorcategoryRepo majorcategoryrepo) Create(majorcategory *model.Majorcategory) (*model.Majorcategory, *httperors.HttpError) {
	if err := majorcategory.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&majorcategory)
	IndexRepo.DbClose(GormDB)
	return majorcategory, nil
}
func (majorcategoryRepo majorcategoryrepo) GetOne(id int) (*model.Majorcategory, *httperors.HttpError) {
	ok := majorcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("majorcategory with that id does not exists!")
	}
	majorcategory := model.Majorcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&majorcategory).Where("id = ?", id).First(&majorcategory)
	IndexRepo.DbClose(GormDB)
	
	return &majorcategory, nil
}

func (majorcategoryRepo majorcategoryrepo) GetAll(majorcategorys []model.Majorcategory,search *support.Search) ([]model.Majorcategory, *httperors.HttpError) {

	results, err1 := majorcategoryRepo.Search(search, majorcategorys)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}
func (majorcategoryRepo majorcategoryrepo) All() (t []model.Majorcategory, r *httperors.HttpError) {

	majorcategory := model.Majorcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	q := GormDB.Model(&majorcategory).Order("name").Find(&t)
	p := paginator.New(adapter.NewGORMAdapter(q), 40)
	p.SetPage(1)

	
	if err3 := p.Results(&t); err3 != nil {
		return nil, httperors.NewNotFoundError("something went wrong paginating!")
	}
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (majorcategoryRepo majorcategoryrepo) Update(id int, majorcategory *model.Majorcategory) (*model.Majorcategory, *httperors.HttpError) {
	ok := majorcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("majorcategory with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	amajorcategory := model.Majorcategory{}
	
	GormDB.Model(&amajorcategory).Where("id = ?", id).First(&amajorcategory)
	if majorcategory.Name  == "" {
		majorcategory.Name = amajorcategory.Name
	}
	if majorcategory.Title  == "" {
		majorcategory.Title = amajorcategory.Title
	}
	if majorcategory.Description  == "" {
		majorcategory.Description = amajorcategory.Description
	}
	fmt.Println(majorcategory)
	GormDB.Model(&amajorcategory).Where("id = ?", id).Update(&majorcategory)
	
	IndexRepo.DbClose(GormDB)

	return majorcategory, nil
}
func (majorcategoryRepo majorcategoryrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := majorcategoryRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	majorcategory := model.Majorcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&majorcategory).Where("id = ?", id).First(&majorcategory)
	GormDB.Delete(majorcategory)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (majorcategoryRepo majorcategoryrepo)ProductUserExistByid(id int) bool {
	majorcategory := model.Majorcategory{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&majorcategory, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (majorcategoryRepo majorcategoryrepo) Search(Ser *support.Search, majorcategorys []model.Majorcategory)([]model.Majorcategory, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	majorcategory := model.Majorcategory{}
	switch(Ser.Search_operator){
	case "all":
		q := GormDB.Model(&majorcategory).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&majorcategorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&majorcategorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&majorcategorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&majorcategorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&majorcategorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&majorcategorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&majorcategorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&majorcategorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&majorcategorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	// break;
	case "like":
		// fmt.Println(Ser.Search_query_1)
		if Ser.Search_query_1 == "all" {
				//db.Order("name DESC")
		q := GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		fmt.Println(p.Results(&majorcategorys))
				if err3 := p.Results(&majorcategorys); err3 != nil {
					return nil, httperors.NewNotFoundError("something went wrong paginating!")
				}

		}else {

			q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);
			p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
			p.SetPage(Ser.Page)
			fmt.Println(p.Results(&majorcategorys,))
			if err3 := p.Results(&majorcategorys); err3 != nil {
				return nil, httperors.NewNotFoundError("something went wrong paginating!")
			}
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&majorcategorys);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&majorcategorys); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return majorcategorys, nil
}