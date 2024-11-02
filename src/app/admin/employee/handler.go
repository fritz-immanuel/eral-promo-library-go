package employee

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/eral-promo-library-go/library"
	"github.com/fritz-immanuel/eral-promo-library-go/middleware"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/library/http/response"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"

	"github.com/fritz-immanuel/eral-promo-library-go/src/services/employee"

	employeeRepository "github.com/fritz-immanuel/eral-promo-library-go/src/services/employee/repository"
	employeeUsecase "github.com/fritz-immanuel/eral-promo-library-go/src/services/employee/usecase"

	employeeroleRepository "github.com/fritz-immanuel/eral-promo-library-go/src/services/employeerole/repository"
	employeeroleUsecase "github.com/fritz-immanuel/eral-promo-library-go/src/services/employeerole/usecase"
)

var ()

type EmployeeHandler struct {
	EmployeeUsecase employee.Usecase
	dataManager     *data.Manager
	Result          gin.H
	Status          int
}

func (h EmployeeHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	employeeRepo := employeeRepository.NewEmployeeRepository(
		data.NewMySQLStorage(db, "employees", models.Employee{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	employeebrandRepo := employeeRepository.NewEmployeeBrandRepository(
		data.NewMySQLStorage(db, "employee_brands", models.EmployeeBrand{}, data.MysqlConfig{}),
	)

	employeeroleRepo := employeeroleRepository.NewEmployeeRoleRepository(
		data.NewMySQLStorage(db, "employee_roles", models.EmployeeRole{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	employeerolepermissionRepo := employeeroleRepository.NewEmployeeRolePermissionRepository(
		data.NewMySQLStorage(db, "employee_role_permissions", models.EmployeeRolePermission{}, data.MysqlConfig{}),
	)

	uEmployeeRole := employeeroleUsecase.NewEmployeeRoleUsecase(db, employeeroleRepo, employeerolepermissionRepo)

	uEmployee := employeeUsecase.NewEmployeeUsecase(db, employeeRepo, employeebrandRepo, uEmployeeRole)

	base := &EmployeeHandler{EmployeeUsecase: uEmployee, dataManager: dataManager}

	rs := v.Group("/employees")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("/:id", middleware.Auth, base.Update)

		rs.PUT("/:id/password", middleware.Auth, base.UpdatePassword)
		rs.PUT("/:id/reset-password", middleware.Auth, base.ResetPassword)

		rs.PUT("/:id/status", middleware.Auth, base.UpdateStatus)
	}

	rss := v.Group("/statuses")
	{
		rss.GET("/employees", base.FindStatus)
	}
}

func (h *EmployeeHandler) FindAll(c *gin.Context) {
	var params models.FindAllEmployeeParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	datas, err := h.EmployeeUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.EmployeeUsecase.Count(c, params)
	if err != nil {
		err.Path = ".EmployeeHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "Employee Data fetched!", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *EmployeeHandler) Find(c *gin.Context) {
	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".EmployeeHandler->Find()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	result, err := h.EmployeeUsecase.Find(c, id)
	if err != nil {
		err.Path = ".EmployeeHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "Employee not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Employee Data fetched!", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	var err *types.Error
	var employee models.Employee
	var dataEmployee *models.Employee

	hash := md5.New()
	io.WriteString(hash, c.PostForm("Password"))

	employee.Name = c.PostForm("Name")
	employee.Email, err = library.IsEmailValid(c.PostForm("Email"))
	if err != nil {
		err.Path = ".EmployeeHandler->Create()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}
	employee.Username = c.PostForm("Username")
	employee.Password = fmt.Sprintf("%x", hash.Sum(nil))
	employee.BusinessID, err = helpers.ValidateUUID(c.PostForm("BusinessID"))
	if err != nil {
		err.Path = ".EmployeeHandler->Create()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}
	employee.EmployeeRoleID, err = helpers.ValidateUUID(c.PostForm("EmployeeRoleID"))
	if err != nil {
		err.Path = ".EmployeeHandler->Create()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	// BRANDS
	errJson := json.Unmarshal([]byte(c.PostForm("Brands")), &employee.Brands)
	if errJson != nil {
		err = &types.Error{
			Path:  ".EmployeeHandler->Create()",
			Error: errJson,
			Type:  "convert-error",
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		dataEmployee, err = h.EmployeeUsecase.Create(c, employee)
		if err != nil {
			return err
		}

		dataEmployee.Password = ""

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".EmployeeHandler->Create()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Employee Data created!", Data: dataEmployee}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	var err *types.Error
	var employee models.Employee
	var data *models.Employee

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".EmployeeHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	employee.Name = c.PostForm("Name")
	employee.Email, err = library.IsEmailValid(c.PostForm("Email"))
	if err != nil {
		err.Path = ".EmployeeHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}
	employee.Username = c.PostForm("Username")
	employee.BusinessID, err = helpers.ValidateUUID(c.PostForm("BusinessID"))
	if err != nil {
		err.Path = ".EmployeeHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}
	employee.EmployeeRoleID, err = helpers.ValidateUUID(c.PostForm("EmployeeRoleID"))
	if err != nil {
		err.Path = ".EmployeeHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	// BRANDS
	errJson := json.Unmarshal([]byte(c.PostForm("Brands")), &employee.Brands)
	if errJson != nil {
		err = &types.Error{
			Path:  ".EmployeeHandler->Update()",
			Error: errJson,
			Type:  "convert-error",
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.EmployeeUsecase.Update(c, id, employee)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".EmployeeHandler->Update()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Employee Data updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeHandler) UpdatePassword(c *gin.Context) {
	var err *types.Error
	var dataEmployee *models.Employee

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".EmployeeHandler->UpdatePassword()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	var oldPassword = c.PostForm("OldPassword")
	var newPassword = c.PostForm("NewPassword")

	if newPassword == "" {
		err = &types.Error{
			Path:    ".EmployeeHandler->UpdatePassword()",
			Message: "New Password cannot be empty",
			Type:    "validation-error",
		}
		response.Error(c, err.Message, http.StatusUnprocessableEntity, *err)
		return
	}

	if newPassword == oldPassword {
		err = &types.Error{
			Path:    ".EmployeeHandler->UpdatePassword()",
			Message: "New Password cannot be the same as the Old Password",
			Type:    "validation-error",
		}
		response.Error(c, err.Message, http.StatusUnprocessableEntity, *err)
		return
	}

	if newPassword != c.PostForm("ConfirmNewPassword") {
		err = &types.Error{
			Path:    ".EmployeeHandler->UpdatePassword()",
			Message: "New Password confirmation does not match",
			Type:    "validation-error",
		}
		response.Error(c, err.Message, http.StatusUnprocessableEntity, *err)
		return
	}

	modelEmployee, err := h.EmployeeUsecase.Find(c, id)
	if err != nil {
		err.Path = ".EmployeeHandler->UpdatePassword()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "Data not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	var currentPassword = modelEmployee.Password
	hash := md5.New()
	io.WriteString(hash, oldPassword)
	hashedOldPassword := fmt.Sprintf("%x", hash.Sum(nil))

	if currentPassword != hashedOldPassword {
		err = &types.Error{
			Path:    ".EmployeeHandler->UpdatePassword()",
			Message: "Incorrect Previous Password",
			Type:    "mysql-error",
		}
		response.Error(c, err.Message, http.StatusUnprocessableEntity, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		dataEmployee, err = h.EmployeeUsecase.UpdatePassword(c, id, newPassword)
		if err != nil {
			return err
		}

		dataEmployee.Password = ""

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".EmployeeHandler->UpdatePassword()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Employee Data updated!", Data: dataEmployee}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeHandler) ResetPassword(c *gin.Context) {
	var err *types.Error
	var dataEmployee *models.Employee

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".EmployeeHandler->ResetPassword()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		dataEmployee, err = h.EmployeeUsecase.UpdatePassword(c, id, "123456")
		if err != nil {
			return err
		}

		dataEmployee.Password = ""

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".EmployeeHandler->ResetPassword()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Employee Password reset successful", Data: dataEmployee}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeHandler) FindStatus(c *gin.Context) {
	var datas []*models.Status
	datas = append(datas, &models.Status{ID: models.STATUS_INACTIVE, Name: "Inactive"})
	datas = append(datas, &models.Status{ID: models.STATUS_ACTIVE, Name: "Active"})

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Employee Status Data fetched!", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.Employee

	employeeID, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".EmployeeHandler->UpdateStatus()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	newStatusID := c.PostForm("StatusID")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.EmployeeUsecase.UpdateStatus(c, employeeID, newStatusID)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".EmployeeHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Employee Status has been updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
