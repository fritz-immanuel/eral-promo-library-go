package employeerole

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/eral-promo-library-go/middleware"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/library/http/response"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"

	"github.com/fritz-immanuel/eral-promo-library-go/src/services/employeerole"

	employeeroleRepository "github.com/fritz-immanuel/eral-promo-library-go/src/services/employeerole/repository"
	employeeroleUsecase "github.com/fritz-immanuel/eral-promo-library-go/src/services/employeerole/usecase"
)

var ()

type EmployeeRoleHandler struct {
	EmployeeRoleUsecase employeerole.Usecase
	dataManager         *data.Manager
	Result              gin.H
	Status              int
}

func (h EmployeeRoleHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	employeeroleRepo := employeeroleRepository.NewEmployeeRoleRepository(
		data.NewMySQLStorage(db, "employee_roles", models.EmployeeRole{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	employeerolepermissionRepo := employeeroleRepository.NewEmployeeRolePermissionRepository(
		data.NewMySQLStorage(db, "employee_role_permissions", models.EmployeeRolePermission{}, data.MysqlConfig{}),
	)

	uEmployeeRole := employeeroleUsecase.NewEmployeeRoleUsecase(db, employeeroleRepo, employeerolepermissionRepo)

	base := &EmployeeRoleHandler{EmployeeRoleUsecase: uEmployeeRole, dataManager: dataManager}

	rs := v.Group("/employee-roles")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("/:id", middleware.Auth, base.Update)

		rs.PUT("/:id/status", middleware.Auth, base.UpdateStatus)
	}

	rss := v.Group("/statuses")
	{
		rss.GET("/employee-roles", base.FindStatus)
	}
}

func (h *EmployeeRoleHandler) FindAll(c *gin.Context) {
	var params models.FindAllEmployeeRoleParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	datas, err := h.EmployeeRoleUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.EmployeeRoleUsecase.Count(c, params)
	if err != nil {
		err.Path = ".EmployeeRoleHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "Employee Role Data fetched!", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *EmployeeRoleHandler) Find(c *gin.Context) {
	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".EmployeeRoleHandler->Find()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	result, err := h.EmployeeRoleUsecase.Find(c, id)
	if err != nil {
		err.Path = ".EmployeeRoleHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "EmployeeRole not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Employee Role Data fetched!", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeRoleHandler) Create(c *gin.Context) {
	var err *types.Error
	var employeerole models.EmployeeRole
	var dataEmployeeRole *models.EmployeeRole

	employeerole.Name = c.PostForm("Name")
	employeerole.IsSupervisor, _ = strconv.Atoi(c.PostForm("IsSupervisor"))

	errJson := json.Unmarshal([]byte(c.PostForm("Permission")), &employeerole.Permission)
	if errJson != nil {
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, types.Error{
			Path:  ".EmployeeRoleHandler->Create()",
			Error: errJson,
			Type:  "convert-error",
		})
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		dataEmployeeRole, err = h.EmployeeRoleUsecase.Create(c, employeerole)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".EmployeeRoleHandler->Create()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Employee Role Data created!", Data: dataEmployeeRole}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeRoleHandler) Update(c *gin.Context) {
	var err *types.Error
	var employeerole models.EmployeeRole
	var data *models.EmployeeRole

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".EmployeeRoleHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	employeerole.Name = c.PostForm("Name")
	employeerole.IsSupervisor, _ = strconv.Atoi(c.PostForm("IsSupervisor"))

	errJson := json.Unmarshal([]byte(c.PostForm("Permission")), &employeerole.Permission)
	if errJson != nil {
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, types.Error{
			Path:  ".EmployeeRoleHandler->Update()",
			Error: errJson,
			Type:  "convert-error",
		})
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.EmployeeRoleUsecase.Update(c, id, employeerole)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".EmployeeRoleHandler->Update()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Employee Role Data updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeRoleHandler) FindStatus(c *gin.Context) {
	var datas []*models.Status
	datas = append(datas, &models.Status{ID: models.STATUS_INACTIVE, Name: "Inactive"})
	datas = append(datas, &models.Status{ID: models.STATUS_ACTIVE, Name: "Active"})

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "EmployeeRole Status Data fetched!", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(http.StatusOK, h.Result)
}

func (h *EmployeeRoleHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.EmployeeRole

	employeeroleID, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".EmployeeRoleHandler->UpdateStatus()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	newStatusID := c.PostForm("StatusID")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.EmployeeRoleUsecase.UpdateStatus(c, employeeroleID, newStatusID)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".EmployeeRoleHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "EmployeeRole Status has been updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
