package employee

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type Repository interface {
	FindAll(*gin.Context, models.FindAllEmployeeParams) ([]*models.Employee, *types.Error)
	Find(*gin.Context, string) (*models.Employee, *types.Error)
	Create(*gin.Context, *models.Employee) (*models.Employee, *types.Error)
	Update(*gin.Context, *models.Employee) (*models.Employee, *types.Error)
	UpdateStatus(*gin.Context, string, string) (*models.Employee, *types.Error)

	FindAllForLogin(*gin.Context, models.FindAllEmployeeParams) ([]*models.EmployeeListForLogin, *types.Error)
}

type BrandRepository interface {
	FindAll(*gin.Context, models.FindAllEmployeeBrandParams) ([]*models.EmployeeBrand, *types.Error)
	Find(*gin.Context, string) (*models.EmployeeBrand, *types.Error)
	Create(*gin.Context, *models.EmployeeBrand) (*models.EmployeeBrand, *types.Error)
	Update(*gin.Context, *models.EmployeeBrand) (*models.EmployeeBrand, *types.Error)
	DeleteByEmployeeID(*gin.Context, string) *types.Error
}
