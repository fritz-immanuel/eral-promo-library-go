package employee

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jmoiron/sqlx"

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
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("", middleware.Auth, base.Update)

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

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Employee Berhasil Ditampilkan", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *EmployeeHandler) Find(c *gin.Context) {
	id := c.Param("id")

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

func (h *EmployeeHandler) Create(c *gin.Context) {
	var err *types.Error
	var employee models.Employee
	var dataEmployee *models.Employee

	hash := md5.New()
	io.WriteString(hash, c.PostForm("Password"))

	employee.Name = c.PostForm("Name")
	employee.Email = c.PostForm("Email")
	employee.Username = c.PostForm("Username")
	employee.Password = fmt.Sprintf("%x", hash.Sum(nil))

	errJson := json.Unmarshal([]byte(c.PostForm("Permission")), &employee.Permission)
	if errJson != nil {
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, types.Error{
			Path:  ".EmployeeHandler->Create()",
			Error: errJson,
			Type:  "convert-error",
		})
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

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Employee Berhasil Ditambahkan", Data: dataEmployee}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	var err *types.Error
	var employee models.Employee
	var data *models.Employee

	id := c.Param("id")

	employee.Name = c.PostForm("Name")
	employee.Email = c.PostForm("Email")
	employee.Username = c.PostForm("Username")

	errJson := json.Unmarshal([]byte(c.PostForm("Permission")), &employee.Permission)
	if errJson != nil {
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, types.Error{
			Path:  ".EmployeeHandler->Update()",
			Error: errJson,
			Type:  "convert-error",
		})
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

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Employee Berhasil Ditambahkan", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeHandler) UpdatePassword(c *gin.Context) {
	var err *types.Error
	var dataEmployee *models.Employee

	id := c.Param("id")
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

func (h *EmployeeHandler) ResetPassword(c *gin.Context) {
	var err *types.Error
	var dataEmployee *models.Employee

	id := c.Param("id")

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

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Employee berhasil direset password", Data: dataEmployee}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeHandler) FindStatus(c *gin.Context) {
	var datas []*models.Status
	datas = append(datas, &models.Status{ID: models.STATUS_INACTIVE, Name: "Inactive"})
	datas = append(datas, &models.Status{ID: models.STATUS_ACTIVE, Name: "Active"})

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Employee Berhasil Ditampilkan", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.Employee

	employeeID := c.Param("id")

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

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Status Employee Berhasil Diubah", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
