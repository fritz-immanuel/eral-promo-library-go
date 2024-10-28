package brand

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/eral-promo-library-go/library/firebase"
	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/middleware"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/brand"
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/http/response"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"

	brandRepository "github.com/fritz-immanuel/eral-promo-library-go/src/services/brand/repository"
	brandUsecase "github.com/fritz-immanuel/eral-promo-library-go/src/services/brand/usecase"
)

type BrandHandler struct {
	BrandUsecase brand.Usecase
	dataManager  *data.Manager
	Result       gin.H
	Status       int
}

func (h BrandHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	brandRepo := brandRepository.NewBrandRepository(
		data.NewMySQLStorage(db, "brands", models.Brand{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	uBrand := brandUsecase.NewBrandUsecase(db, brandRepo)

	base := &BrandHandler{BrandUsecase: uBrand, dataManager: dataManager}

	rs := v.Group("/brands")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("/:id", middleware.Auth, base.Update)

		rs.PUT("/:id/status", middleware.Auth, base.UpdateStatus)
	}

	rss := v.Group("/statuses")
	{
		rss.GET("/brands", base.FindStatus)
	}
}

func (h *BrandHandler) FindAll(c *gin.Context) {
	var params models.FindAllBrandParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	params.FindAllParams.SortBy = "brands.name ASC"
	datas, err := h.BrandUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}
	}

	for _, data := range datas {
		if data.LogoImgURL != "" {
			data.LogoImgURL, _ = firebase.GenerateSignedURL(data.LogoImgURL)
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.BrandUsecase.Count(c, params)
	if err != nil {
		err.Path = ".BrandHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "Brand Data fetched!", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *BrandHandler) Find(c *gin.Context) {
	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".BrandHandler->Find()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	result, err := h.BrandUsecase.Find(c, id)
	if err != nil {
		err.Path = ".BrandHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "Brand not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	if result.LogoImgURL != "" {
		result.LogoImgURL, _ = firebase.GenerateSignedURL(result.LogoImgURL)
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Brand Data fetched!", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *BrandHandler) Create(c *gin.Context) {
	var err *types.Error
	var brand models.Brand
	var dataBrand *models.Brand

	brand.Name = c.PostForm("Name")
	brand.Code = c.PostForm("Code")
	brand.BusinessID = c.PostForm("BusinessID")

	file, errFile := c.FormFile("LogoImgURL")
	if file != nil {
		if errFile != nil {
			err = &types.Error{
				Path:       ".BrandHandler->Create()",
				Message:    errFile.Error(),
				Error:      errFile,
				StatusCode: http.StatusInternalServerError,
				Type:       "golang-error",
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		filename, err := firebase.UploadFile(c, file, "brand")
		if err != nil {
			err.Path = ".BrandHandler->Create()" + err.Path
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		brand.LogoImgURL = filename
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		dataBrand, err = h.BrandUsecase.Create(c, brand)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".BrandHandler->Create()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Brand Data created!", Data: dataBrand}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *BrandHandler) Update(c *gin.Context) {
	var err *types.Error
	var brand models.Brand
	var data *models.Brand

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".BrandHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	brand.Name = c.PostForm("Name")
	brand.Code = c.PostForm("Code")
	brand.BusinessID = c.PostForm("BusinessID")

	file, errFile := c.FormFile("LogoImgURL")
	if file != nil {
		if errFile != nil {
			err = &types.Error{
				Path:       ".BrandHandler->Update()",
				Message:    errFile.Error(),
				Error:      errFile,
				StatusCode: http.StatusInternalServerError,
				Type:       "golang-error",
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		filename, err := firebase.UploadFile(c, file, "brand")
		if err != nil {
			err.Path = ".BrandHandler->Update()" + err.Path
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		brand.LogoImgURL = filename
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		{ // delete existing logo
			brandData, err := h.BrandUsecase.Find(c, id)
			if err != nil {
				return err
			}

			if brandData.LogoImgURL != "" {
				err := firebase.DeleteFile(c, brandData.LogoImgURL)
				if err != nil {
					return err
				}
			}
		}

		data, err = h.BrandUsecase.Update(c, id, brand)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".BrandHandler->Update()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Brand Data updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *BrandHandler) FindStatus(c *gin.Context) {
	var datas []*models.Status
	datas = append(datas, &models.Status{ID: models.STATUS_INACTIVE, Name: "Inactive"})
	datas = append(datas, &models.Status{ID: models.STATUS_ACTIVE, Name: "Active"})

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Brand Status Data fetched!", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(http.StatusOK, h.Result)
}

func (h *BrandHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.Brand

	brandID, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".BrandHandler->UpdateStatus()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	newStatusID := c.PostForm("StatusID")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.BrandUsecase.UpdateStatus(c, brandID, newStatusID)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".BrandHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Brand Status has been updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
