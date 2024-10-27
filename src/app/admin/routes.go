package admin

import (
	http_business "github.com/fritz-immanuel/eral-promo-library-go/src/app/admin/business"
	http_businessconfig "github.com/fritz-immanuel/eral-promo-library-go/src/app/admin/businessconfig"
	http_company "github.com/fritz-immanuel/eral-promo-library-go/src/app/admin/company"
	http_employee "github.com/fritz-immanuel/eral-promo-library-go/src/app/admin/employee"
	http_employeerole "github.com/fritz-immanuel/eral-promo-library-go/src/app/admin/employeerole"
	http_user "github.com/fritz-immanuel/eral-promo-library-go/src/app/admin/user"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var (
	businessHandler       http_business.BusinessHandler
	businessconfigHandler http_businessconfig.BusinessConfigHandler
	companyHandler        http_company.CompanyHandler
	employeeHandler       http_employee.EmployeeHandler
	employeeroleHandler   http_employeerole.EmployeeRoleHandler
	userHandler           http_user.UserHandler
)

func RegisterRoutes(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	v1 := v.Group("")
	{
		businessHandler.RegisterAPI(db, dataManager, router, v1)
		businessconfigHandler.RegisterAPI(db, dataManager, router, v1)
		companyHandler.RegisterAPI(db, dataManager, router, v1)
		employeeHandler.RegisterAPI(db, dataManager, router, v1)
		employeeroleHandler.RegisterAPI(db, dataManager, router, v1)
		userHandler.RegisterAPI(db, dataManager, router, v1)
	}
}
