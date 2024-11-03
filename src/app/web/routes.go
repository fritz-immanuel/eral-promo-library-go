package web

import (
	http_brand "github.com/fritz-immanuel/eral-promo-library-go/src/app/web/brand"
	http_business "github.com/fritz-immanuel/eral-promo-library-go/src/app/web/business"
	http_businessconfig "github.com/fritz-immanuel/eral-promo-library-go/src/app/web/businessconfig"
	http_company "github.com/fritz-immanuel/eral-promo-library-go/src/app/web/company"
	http_employee "github.com/fritz-immanuel/eral-promo-library-go/src/app/web/employee"
	http_promo "github.com/fritz-immanuel/eral-promo-library-go/src/app/web/promo"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var (
	brandHandler          http_brand.BrandHandler
	businessHandler       http_business.BusinessHandler
	businessconfigHandler http_businessconfig.BusinessConfigHandler
	companyHandler        http_company.CompanyHandler
	employeeHandler       http_employee.EmployeeHandler
	promoHandler          http_promo.PromoHandler
)

func RegisterRoutes(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	v1 := v.Group("")
	{
		brandHandler.RegisterAPI(db, dataManager, router, v1)
		businessHandler.RegisterAPI(db, dataManager, router, v1)
		businessconfigHandler.RegisterAPI(db, dataManager, router, v1)
		companyHandler.RegisterAPI(db, dataManager, router, v1)
		employeeHandler.RegisterAPI(db, dataManager, router, v1)
		promoHandler.RegisterAPI(db, dataManager, router, v1)
	}
}
