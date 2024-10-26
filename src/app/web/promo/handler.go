package promo

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/eral-promo-library-go/library"
	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/middleware"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/promo"
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/http/response"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"

	promoRepository "github.com/fritz-immanuel/eral-promo-library-go/src/services/promo/repository"
	promoUsecase "github.com/fritz-immanuel/eral-promo-library-go/src/services/promo/usecase"
)

type PromoHandler struct {
	PromoUsecase promo.Usecase
	dataManager  *data.Manager
	Result       gin.H
	Status       int
}

func (h PromoHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	promoRepo := promoRepository.NewPromoRepository(
		data.NewMySQLStorage(db, "promos", models.Promo{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "promo_status", models.Status{}, data.MysqlConfig{}),
	)

	promodocumentRepo := promoRepository.NewPromoDocumentRepository(
		data.NewMySQLStorage(db, "promo_documents", models.PromoDocument{}, data.MysqlConfig{}),
	)

	uPromo := promoUsecase.NewPromoUsecase(db, promoRepo, promodocumentRepo)

	base := &PromoHandler{PromoUsecase: uPromo, dataManager: dataManager}

	rs := v.Group("/promos")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("", middleware.Auth, base.Update)

		rs.PUT("/:id/status", middleware.Auth, base.UpdateStatus)
	}

	rss := v.Group("/statuses")
	{
		rss.GET("/promos", base.FindStatus)
	}
}

func (h *PromoHandler) FindAll(c *gin.Context) {
	var params models.FindAllPromoParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	params.FindAllParams.SortBy = "promos.name ASC"
	datas, err := h.PromoUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.PromoUsecase.Count(c, params)
	if err != nil {
		err.Path = ".PromoHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Promo Berhasil Ditampilkan", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *PromoHandler) Find(c *gin.Context) {
	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".PromoHandler->Find()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	result, err := h.PromoUsecase.Find(c, *id)
	if err != nil {
		err.Path = ".PromoHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "Promo not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Promo Berhasil Ditampilkan", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *PromoHandler) Create(c *gin.Context) {
	var err *types.Error
	var promo models.Promo
	var dataPromo *models.Promo

	if c.Query("StartDate") != "" {
		startDateTime, errConversion := time.Parse(library.DateStampFormat(), c.Query("StartDate"))
		if errConversion != nil {
			err := &types.Error{
				Path:       ".PromoHandler->Create()",
				Message:    "Incorrect Start Date Format",
				Error:      errConversion,
				Type:       "conversion-error",
				StatusCode: http.StatusBadRequest,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}
		promo.StartDate = startDateTime
	}

	if c.Query("EndDate") != "" {
		endDateTime, errConversion := time.Parse(library.DateStampFormat(), c.Query("EndDate"))
		if errConversion != nil {
			err := &types.Error{
				Path:       ".PromoHandler->Create()",
				Message:    "Incorrect End Date Format",
				Error:      errConversion,
				Type:       "conversion-error",
				StatusCode: http.StatusBadRequest,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}
		promo.EndDate = endDateTime
	}

	promo.Name = c.PostForm("Name")
	promo.Code = c.PostForm("Code")
	promo.PromoTypeID = c.PostForm("PromoTypeID")
	promo.CompanyID = c.PostForm("CompanyID")
	promo.BusinessID = c.PostForm("BusinessID")
	promo.TotalPromoBudget, _ = strconv.ParseFloat(c.PostForm("TotalPromoBudget"), 64)
	promo.PrincipleSupport, _ = strconv.ParseFloat(c.PostForm("PrincipleSupport"), 64)
	promo.InternalSupport, _ = strconv.ParseFloat(c.PostForm("InternalSupport"), 64)
	promo.Description = c.PostForm("Description")

	// TODO: upload img
	// TODO: upload documents

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		dataPromo, err = h.PromoUsecase.Create(c, promo)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".PromoHandler->Create()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Promo Berhasil Ditambahkan", Data: dataPromo}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *PromoHandler) Update(c *gin.Context) {
	var err *types.Error
	var promo models.Promo
	var data *models.Promo

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".PromoHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if c.Query("StartDate") != "" {
		startDateTime, errConversion := time.Parse(library.DateStampFormat(), c.Query("StartDate"))
		if errConversion != nil {
			err := &types.Error{
				Path:       ".PromoHandler->Update()",
				Message:    "Incorrect Start Date Format",
				Error:      errConversion,
				Type:       "conversion-error",
				StatusCode: http.StatusBadRequest,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}
		promo.StartDate = startDateTime
	}

	if c.Query("EndDate") != "" {
		endDateTime, errConversion := time.Parse(library.DateStampFormat(), c.Query("EndDate"))
		if errConversion != nil {
			err := &types.Error{
				Path:       ".PromoHandler->Update()",
				Message:    "Incorrect End Date Format",
				Error:      errConversion,
				Type:       "conversion-error",
				StatusCode: http.StatusBadRequest,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}
		promo.EndDate = endDateTime
	}

	promo.Name = c.PostForm("Name")
	promo.Code = c.PostForm("Code")
	promo.PromoTypeID = c.PostForm("PromoTypeID")
	promo.CompanyID = c.PostForm("CompanyID")
	promo.BusinessID = c.PostForm("BusinessID")
	promo.TotalPromoBudget, _ = strconv.ParseFloat(c.PostForm("TotalPromoBudget"), 64)
	promo.PrincipleSupport, _ = strconv.ParseFloat(c.PostForm("PrincipleSupport"), 64)
	promo.InternalSupport, _ = strconv.ParseFloat(c.PostForm("InternalSupport"), 64)
	promo.Description = c.PostForm("Description")

	// TODO: upload img
	// TODO: upload documents

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.PromoUsecase.Update(c, *id, promo)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".PromoHandler->Update()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Promo Berhasil Ditambahkan", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *PromoHandler) FindStatus(c *gin.Context) {
	datas, err := h.PromoUsecase.FindStatus(c)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Promo Status fetched!", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(http.StatusOK, h.Result)
}

func (h *PromoHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.Promo

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".PromoHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	newStatusID := c.PostForm("StatusID")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.PromoUsecase.UpdateStatus(c, *id, newStatusID)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".PromoHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Promo Status has been updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
