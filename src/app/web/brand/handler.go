package brand

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/eral-promo-library-go/library/appcontext"
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
		rs.GET("", middleware.AuthWebApp, base.FindAll)
		rs.GET("/:id", middleware.AuthWebApp, base.Find)
	}
}

func (h *BrandHandler) FindAll(c *gin.Context) {
	var params models.FindAllBrandParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	params.BusinessID = *appcontext.BusinessID(c)
	params.FindAllParams.DataFinder = fmt.Sprintf(`brands.id IN (SELECT brand_id FROM employee_brands WHERE employee_id = "%s")`, *appcontext.EmployeeID(c))
	datas, err := h.BrandUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
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
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
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
