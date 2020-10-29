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
	InvoiceController invoiceController = invoiceController{}
)
type invoiceController struct{ }
/////////controllers/////////////////
func (controller invoiceController) Create(c echo.Context) error {
	invoice := &model.Invoice{}
	
	if err := c.Bind(invoice); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	createdinvoice, err1 := service.Invoiceservice.Create(invoice)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdinvoice)
}
func (controller invoiceController) GetAll(c echo.Context) error {
	invoices := []model.Invoice{}
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
	
	invoices, err3 := service.Invoiceservice.GetAll(invoices,search)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, invoices)
}
 func (controller invoiceController) View(c echo.Context) error {
	code, problem := service.Invoiceservice.View()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, code)	
}
func (controller invoiceController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	invoice, problem := service.Invoiceservice.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, invoice)	
}

func (controller invoiceController) Update(c echo.Context) error {
		
	invoice :=  &model.Invoice{}
	if err := c.Bind(invoice); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedinvoice, problem := service.Invoiceservice.Update(id, invoice)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedinvoice)
}

func (controller invoiceController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.Invoiceservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
