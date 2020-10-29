package repository

import (
	// "log"
	// "os"
	// "github.com/joho/godotenv"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
	// _ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/myrachanto/astore-v2.0/httperors" 
	"github.com/myrachanto/astore-v2.0/model"
)
var (
	IndexRepo indexRepo = indexRepo{}
	Operator = map[string]string{"all":"all","equal_to":"=","not_equal_to":"<>","less_than":"<",
"greater_than":">","less_than_or_equal_to":"<=","greater_than_ro_equal_to":">=",
"like":"like","between":"between","in":"in","not_in":"not_in"}
)

///curtesy to gorm
type indexRepo struct{}
func (indexRepo indexRepo) Getconnected() (GormDB *gorm.DB, err *httperors.HttpError) {
	// err1 := godotenv.Load()
	// if err1 != nil {
	// 	log.Fatal("Error loading .env file in routes")
	// }
	// dbuser := os.Getenv("DbUsername")
	// DbName := os.Getenv("DbName")
	// dbURI := dbuser+"@/"+DbName+"?charset=utf8&parseTime=True&loc=Local"
	// GormDB, err2 := gorm.Open("mysql", dbURI)
	// if err2 != nil {
	// 	return nil, httperors.NewNotFoundError("No Mysql db connection")
	// }
	GormDB, err1 := gorm.Open("sqlite3", "estore.db")
    if err1 != nil {
        panic("failed to connect database")
    }
	GormDB.AutoMigrate(&model.Cart{})
	GormDB.AutoMigrate(&model.Category{})
	GormDB.AutoMigrate(&model.Customer{})
	GormDB.AutoMigrate(&model.Invoice{})
	GormDB.AutoMigrate(&model.Majorcategory{})
	GormDB.AutoMigrate(&model.Product{})
	GormDB.AutoMigrate(&model.Subcategory{})
	GormDB.AutoMigrate(&model.Transaction{})
	GormDB.AutoMigrate(&model.User{})
	GormDB.AutoMigrate(&model.Auth{})
	return GormDB, nil
}
func (indexRepo indexRepo) DbClose(GormDB *gorm.DB) {
	defer GormDB.Close()
}