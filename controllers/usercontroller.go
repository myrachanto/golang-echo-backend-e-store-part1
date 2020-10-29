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
	UserController userController = userController{}
)
type userController struct{ }
/////////controllers/////////////////
func (controller userController) Create(c echo.Context) error {
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		httperror := httperors.NewBadRequestError("Invalid user json body")
		return c.JSON(httperror.Code, httperror)
	}
	s, err1 := service.UserService.Create(user)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, s)
}
func (controller userController) Login(c echo.Context) error {
	user := &model.LoginUser{}
	auth := &model.Auth{}
	if err := c.Bind(user); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	
	auth, problem := service.UserService.Login(user)
	if problem != nil {
		fmt.Println(problem)
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, auth)	
}
func (controller userController) Logout(c echo.Context) error {
	token := string(c.Param("token"))
	problem := service.UserService.Logout(token)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, "succeessifully logged out")	
}

func (controller userController) GetAll(c echo.Context) error {
	Users := []model.User{}
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
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid per number")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println("------------------------")
	search := &support.Search{Column:column, Direction:direction,Search_column:search_column,Search_operator:search_operator,Search_query_1:search_query_1,Search_query_2:search_query_2,Per_page:per_page,Page:page}
	
	users, err3 := service.UserService.GetAll(Users,search)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, users)
} 
func (controller userController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	user, problem := service.UserService.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, user)	
}

func (controller userController) Update(c echo.Context) error {
	user :=  &model.User{}
	if err := c.Bind(user); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updateduser, problem := service.UserService.Update(id, user)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updateduser)
}

func (controller userController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.UserService.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}