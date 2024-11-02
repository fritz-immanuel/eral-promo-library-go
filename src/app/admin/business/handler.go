package business

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/eral-promo-library-go/library/firebase"
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
		rs.PUT("/:id", middleware.Auth, base.Update)

		rs.PUT("/:id/status", middleware.Auth, base.UpdateStatus)
	}

	rss := v.Group("/statuses")
	{
		rss.GET("/business", base.FindStatus)
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
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}
	}

	for _, data := range datas {
		if data.LogoImgURL != "" {
			data.LogoImgURL, err = firebase.GenerateSignedURL(data.LogoImgURL)
			if err != nil {
				fmt.Println(err.Error)
			}
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

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "Business Data fetched!", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *BusinessHandler) Find(c *gin.Context) {
	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".BusinessHandler->Find()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

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

	if result.LogoImgURL != "" {
		result.LogoImgURL, _ = firebase.GenerateSignedURL(result.LogoImgURL)
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Business Data fetched!", Data: result}
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
	business.CompanyID = c.PostForm("CompanyID")

	file, errFile := c.FormFile("LogoImgURL")
	if file != nil {
		if errFile != nil {
			err = &types.Error{
				Path:       ".BusinessHandler->Create()",
				Message:    errFile.Error(),
				Error:      errFile,
				StatusCode: http.StatusInternalServerError,
				Type:       "golang-error",
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		filename, err := firebase.UploadFile(c, file, "business")
		if err != nil {
			err.Path = ".BusinessHandler->Create()" + err.Path
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		business.LogoImgURL = filename
	}

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

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Business Data created!", Data: dataBusiness}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *BusinessHandler) Update(c *gin.Context) {
	var err *types.Error
	var business models.Business
	var data *models.Business

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".BusinessHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	business.Name = c.PostForm("Name")
	business.Code = c.PostForm("Code")
	business.CompanyID = c.PostForm("CompanyID")

	file, errFile := c.FormFile("LogoImgURL")
	if file != nil {
		if errFile != nil {
			err = &types.Error{
				Path:       ".BusinessHandler->Update()",
				Message:    errFile.Error(),
				Error:      errFile,
				StatusCode: http.StatusInternalServerError,
				Type:       "golang-error",
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		filename, err := firebase.UploadFile(c, file, "business")
		if err != nil {
			err.Path = ".BusinessHandler->Update()" + err.Path
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		business.LogoImgURL = filename
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		{ // delete existing logo
			businessData, err := h.BusinessUsecase.Find(c, id)
			if err != nil {
				return err
			}

			if businessData.LogoImgURL != "" {
				err := firebase.DeleteFile(c, businessData.LogoImgURL)
				if err != nil {
					return err
				}
			}
		}

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

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Business Data updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *BusinessHandler) FindStatus(c *gin.Context) {
	var datas []*models.Status
	datas = append(datas, &models.Status{ID: models.STATUS_INACTIVE, Name: "Inactive"})
	datas = append(datas, &models.Status{ID: models.STATUS_ACTIVE, Name: "Active"})

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Business Status Data fetched!", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(http.StatusOK, h.Result)
}

func (h *BusinessHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.Business

	businessID, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".BusinessHandler->UpdateStatus()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	newStatusID := c.PostForm("StatusID")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.BusinessUsecase.UpdateStatus(c, businessID, newStatusID)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".BusinessHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Business Status has been updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
