package company

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/eral-promo-library-go/library/appcontext"
	"github.com/fritz-immanuel/eral-promo-library-go/library/firebase"
	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/middleware"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/company"
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/http/response"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"

	companyRepository "github.com/fritz-immanuel/eral-promo-library-go/src/services/company/repository"
	companyUsecase "github.com/fritz-immanuel/eral-promo-library-go/src/services/company/usecase"
)

type CompanyHandler struct {
	CompanyUsecase company.Usecase
	dataManager    *data.Manager
	Result         gin.H
	Status         int
}

func (h CompanyHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	companyRepo := companyRepository.NewCompanyRepository(
		data.NewMySQLStorage(db, "companies", models.Company{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	uCompany := companyUsecase.NewCompanyUsecase(db, companyRepo)

	base := &CompanyHandler{CompanyUsecase: uCompany, dataManager: dataManager}

	rs := v.Group("/company")
	{
		rs.GET("", middleware.AuthWebApp, base.Find)
	}
}

func (h *CompanyHandler) Find(c *gin.Context) {
	id, err := helpers.ValidateUUID(*appcontext.CompanyID(c))
	if err != nil {
		err.Path = ".CompanyHandler->Find()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	result, err := h.CompanyUsecase.Find(c, id)
	if err != nil {
		err.Path = ".CompanyHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "Company not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	if result.LogoImgURL != "" {
		result.LogoImgURL, _ = firebase.GenerateSignedURL(result.LogoImgURL)
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Company Data fetched!", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
