package business

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/eral-promo-library-go/library/appcontext"
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
		rs.GET("", middleware.AuthWebApp, base.Find)
	}
}

func (h *BusinessHandler) Find(c *gin.Context) {
	id, err := helpers.ValidateUUID(*appcontext.BusinessID(c))
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
