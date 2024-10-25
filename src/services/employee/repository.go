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
}

type PermissionRepository interface {
	FindAll(*gin.Context, models.FindAllEmployeePermissionParams) ([]*models.EmployeePermission, *types.Error)
	Find(*gin.Context, string) (*models.EmployeePermission, *types.Error)
	Create(*gin.Context, *models.CreateUpdateEmployeePermission) (*models.EmployeePermission, *types.Error)
	DeleteByEmployeeID(*gin.Context, string) *types.Error

	CreateBunch(*gin.Context, string, models.FindAllEmployeePermissionParams) *types.Error
}
