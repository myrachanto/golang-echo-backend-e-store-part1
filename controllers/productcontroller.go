package controllers

import(
	"fmt"
	"strconv"	
	"net/http"
	"github.com/labstack/echo"
	"github.com/myrachanto/astore-v2.0/httperors"
	"github.com/myrachanto/astore-v2.0/model"
	"github.com/myrachanto/astore-v2.0/service"
	"github.com/myrachanto/astore-v2.0/support"
)
 
var (
	ProductController productController = productController{}
)
type productController struct{ }
/////////controllers/////////////////
func (controller productController) Create(c echo.Context) error {
	product := &model.Product{}
	
	if err := c.Bind(product); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	createdproduct, err1 := service.Productservice.Create(product)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	} 
	return c.JSON(http.StatusCreated, createdproduct)
}

func (controller productController) View(c echo.Context) error {
	options, problem := service.Productservice.View()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}
func (controller productController) GetProducts(c echo.Context) error {
	products := []model.Product{}	
	column := string(c.QueryParam("column"))
	name := string(c.QueryParam("name"))
	search_column := string(c.QueryParam("search_column"))
	search_operator := string(c.QueryParam("search_operator"))
	search_query_1 := string(c.QueryParam("search_query_1"))
	per_page, err := strconv.Atoi(c.QueryParam("per_page"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid per number")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println("------------------------")
	search := &support.Productsearch{Column:column,Name:name,Search_column:search_column,Search_operator:search_operator,Search_query_1:search_query_1,Per_page:per_page}
	
	products, err3 := service.Productservice.GetProducts(products, search)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, products)
} 
func (controller productController) GetAll(c echo.Context) error {
	products := []model.Product{}	
	column := string(c.QueryParam("column"))
	direction := string(c.QueryParam("direction"))
	search_column := string(c.QueryParam("search_column"))
	search_operator := string(c.QueryParam("search_operator"))
	search_query_1 := string(c.QueryParam("search_query_1"))
	search_query_2 := string(c.QueryParam("search_query_2"))
	per_page, err := strconv.Atoi(c.QueryParam("per_page"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid per number")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println("------------------------")
	search := &support.Search{Column:column, Direction:direction,Search_column:search_column,Search_operator:search_operator,Search_query_1:search_query_1,Search_query_2:search_query_2,Per_page:per_page}
	
	products, err3 := service.Productservice.GetAll(products, search)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, products)
} 
func (controller productController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	product, problem := service.Productservice.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, product)	
}

func (controller productController) Update(c echo.Context) error {
		
	product :=  &model.Product{}
	
	if err := c.Bind(product); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedproduct, problem := service.Productservice.Update(id, product)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedproduct)
}

func (controller productController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.Productservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
