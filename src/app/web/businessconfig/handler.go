package businessconfig

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/businessconfig"
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/http/response"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"

	businessconfigRepository "github.com/fritz-immanuel/eral-promo-library-go/src/services/businessconfig/repository"
	businessconfigUsecase "github.com/fritz-immanuel/eral-promo-library-go/src/services/businessconfig/usecase"
)

type BusinessConfigHandler struct {
	BusinessConfigUsecase businessconfig.Usecase
	dataManager           *data.Manager
	Result                gin.H
	Status                int
}

func (h BusinessConfigHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	businessconfigRepo := businessconfigRepository.NewBusinessConfigRepository(
		data.NewMySQLStorage(db, "business_configs", models.BusinessConfig{}, data.MysqlConfig{}),
	)

	uBusinessConfig := businessconfigUsecase.NewBusinessConfigUsecase(db, businessconfigRepo)

	base := &BusinessConfigHandler{BusinessConfigUsecase: uBusinessConfig, dataManager: dataManager}

	rs := v.Group("/config/business")
	{
		rs.GET("", base.FindAll)
	}
}

func (h *BusinessConfigHandler) FindAll(c *gin.Context) {
	var params models.FindAllBusinessConfigParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	params.FindAllParams.SortBy = "business_configs.id ASC"
	params.SubURLName = c.Query("SubURLName")
	datas, err := h.BusinessConfigUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.BusinessConfigUsecase.Count(c, params)
	if err != nil {
		err.Path = ".BusinessConfigHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "Business Config Data fetched!", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}
