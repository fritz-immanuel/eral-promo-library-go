package business

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/middleware"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/business"
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/http/response"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"

	businessRepository "github.com/fritz-immanuel/eral-promo-library-go/src/services/business/repository"
	businessUsecase "github.com/fritz-immanuel/eral-promo-library-go/src/services/business/usecase"
)

type BusinessHandler struct {
	BusinessUsecase business.Usecase
	dataManager     *data.Manager
	Result          gin.H
	Status          int
}

func (h BusinessHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	businessRepo := businessRepository.NewBusinessRepository(
		data.NewMySQLStorage(db, "business", models.Business{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	uBusiness := businessUsecase.NewBusinessUsecase(db, businessRepo)

	base := &BusinessHandler{BusinessUsecase: uBusiness, dataManager: dataManager}

	rs := v.Group("/business")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("", middleware.Auth, base.Update)
	}
}

func (h *BusinessHandler) FindAll(c *gin.Context) {
	var params models.FindAllBusinessParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	params.FindAllParams.SortBy = "business.name ASC"
	datas, err := h.BusinessUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.BusinessUsecase.Count(c, params)
	if err != nil {
		err.Path = ".BusinessHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Business Berhasil Ditampilkan", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *BusinessHandler) Find(c *gin.Context) {
	id := c.Param("id")

	result, err := h.BusinessUsecase.Find(c, id)
	if err != nil {
		err.Path = ".BusinessHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "Business not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Business Berhasil Ditampilkan", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *BusinessHandler) Create(c *gin.Context) {
	var err *types.Error
	var business models.Business
	var dataBusiness *models.Business

	business.Name = c.PostForm("Name")
	business.Code = c.PostForm("Code")

	// TODO: upload img

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		dataBusiness, err = h.BusinessUsecase.Create(c, business)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".BusinessHandler->Create()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Business Berhasil Ditambahkan", Data: dataBusiness}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *BusinessHandler) Update(c *gin.Context) {
	var err *types.Error
	var business models.Business
	var data *models.Business

	id := c.Param("id")

	business.Name = c.PostForm("Name")
	business.Code = c.PostForm("Code")

	// TODO: upload img

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.BusinessUsecase.Update(c, id, business)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".BusinessHandler->Update()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Business Berhasil Ditambahkan", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
