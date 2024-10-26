package employee

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/eral-promo-library-go/middleware"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/eral-promo-library-go/library/appcontext"
	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/http/response"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"

	"github.com/fritz-immanuel/eral-promo-library-go/src/services/employee"

	employeeRepository "github.com/fritz-immanuel/eral-promo-library-go/src/services/employee/repository"
	employeeUsecase "github.com/fritz-immanuel/eral-promo-library-go/src/services/employee/usecase"
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

	employeepermissionRepo := employeeRepository.NewEmployeePermissionRepository(
		data.NewMySQLStorage(db, "employee_permissions", models.EmployeePermission{}, data.MysqlConfig{}),
	)

	uEmployee := employeeUsecase.NewEmployeeUsecase(db, employeeRepo, employeepermissionRepo)

	base := &EmployeeHandler{EmployeeUsecase: uEmployee, dataManager: dataManager}

	rs := v.Group("/employees")
	{
		rs.GET("/profile", middleware.Auth, base.Find)

		rs.PUT("/profile/password", middleware.Auth, base.UpdatePassword)
	}

	rsa := v.Group("/employees/auth")
	{
		rsa.POST("/login", base.Login)
	}
}

func (h *EmployeeHandler) Find(c *gin.Context) {
	id := *appcontext.EmployeeID(c)

	result, err := h.EmployeeUsecase.Find(c, id)
	if err != nil {
		err.Path = ".EmployeeHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "Employee not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Employee Berhasil Ditampilkan", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeHandler) UpdatePassword(c *gin.Context) {
	var err *types.Error
	var dataEmployee *models.Employee

	id := *appcontext.EmployeeID(c)
	var oldPassword = c.PostForm("OldPassword")
	var newPassword = c.PostForm("NewPassword")

	if newPassword == "" {
		err = &types.Error{
			Path:    ".EmployeeHandler->UpdatePassword()",
			Message: "Password baru tidak boleh kosong",
			Type:    "mysql-error",
		}
		response.Error(c, err.Message, http.StatusUnprocessableEntity, *err)
		return
	}

	if newPassword == oldPassword {
		err = &types.Error{
			Path:    ".EmployeeHandler->UpdatePassword()",
			Message: "Password baru tidak boleh sama dengan password lama",
			Type:    "mysql-error",
		}
		response.Error(c, err.Message, http.StatusUnprocessableEntity, *err)
		return
	}

	if newPassword != c.PostForm("ConfirmNewPassword") {
		err = &types.Error{
			Path:    ".EmployeeHandler->UpdatePassword()",
			Message: "Gagal mengkonfirmasi password baru",
			Type:    "mysql-error",
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

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Employee Berhasil Diperbarui", Data: dataEmployee}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeHandler) Login(c *gin.Context) {
	var err *types.Error
	var obj models.EmployeeLogin
	var data *models.EmployeeLogin

	hash := md5.New()
	io.WriteString(hash, c.PostForm("Password"))

	username := c.PostForm("Username")
	password := fmt.Sprintf("%x", hash.Sum(nil))

	obj.Username = username
	obj.Password = password

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.EmployeeUsecase.Login(c, obj)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".EmployeeHandler->Login()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Login berhasil", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
