package web

import (
	http_business "github.com/fritz-immanuel/eral-promo-library-go/src/app/web/business"
	http_businessconfig "github.com/fritz-immanuel/eral-promo-library-go/src/app/web/businessconfig"
	http_employee "github.com/fritz-immanuel/eral-promo-library-go/src/app/web/employee"
	http_promo "github.com/fritz-immanuel/eral-promo-library-go/src/app/web/promo"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var (
	businessHandler       http_business.BusinessHandler
	businessconfigHandler http_businessconfig.BusinessConfigHandler
	employeeHandler       http_employee.EmployeeHandler
	promoHandler          http_promo.PromoHandler
)

func RegisterRoutes(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	v1 := v.Group("")
	{
		businessHandler.RegisterAPI(db, dataManager, router, v1)
		businessconfigHandler.RegisterAPI(db, dataManager, router, v1)
		employeeHandler.RegisterAPI(db, dataManager, router, v1)
		promoHandler.RegisterAPI(db, dataManager, router, v1)
	}
}
