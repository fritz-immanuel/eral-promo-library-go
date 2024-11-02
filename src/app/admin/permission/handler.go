package permission

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/middleware"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/permission"
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/http/response"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"

	permissionRepository "github.com/fritz-immanuel/eral-promo-library-go/src/services/permission/repository"
	permissionUsecase "github.com/fritz-immanuel/eral-promo-library-go/src/services/permission/usecase"
)

var ()

// PermissionHandler  represent the httphandler for article
type PermissionHandler struct {
	PermissionUsecase permission.Usecase
	dataManager       *data.Manager
	Result            gin.H
	Status            int
}

func (h PermissionHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	permissionRepo := permissionRepository.NewPermissionRepository(
		data.NewMySQLStorage(db, "permissions", models.Permission{}, data.MysqlConfig{}),
	)

	uPermission := permissionUsecase.NewPermissionUsecase(db, &permissionRepo)

	base := &PermissionHandler{PermissionUsecase: uPermission, dataManager: dataManager}

	rs := v.Group("/permissions")
	{
		rs.GET("", middleware.Auth, base.FindAll)
	}
}

func (h *PermissionHandler) FindAll(c *gin.Context) {
	var params models.FindAllPermissionParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	params.Package = c.Query("Package")
	params.Name = c.Query("Name")
	params.FindAllParams.DataFinder = "is_hidden = 0"
	params.FindAllParams.SortBy = "module_name ASC, sequence_number_detail ASC"
	datas, err := h.PermissionUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.PermissionUsecase.Count(c, params)
	if err != nil {
		err.Path = ".PermissionHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "Permission Data fetched!", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}
