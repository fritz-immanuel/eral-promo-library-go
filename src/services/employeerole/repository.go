package employeerole

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type Repository interface {
	FindAll(*gin.Context, models.FindAllEmployeeRoleParams) ([]*models.EmployeeRole, *types.Error)
	Find(*gin.Context, string) (*models.EmployeeRole, *types.Error)
	Create(*gin.Context, *models.EmployeeRole) (*models.EmployeeRole, *types.Error)
	Update(*gin.Context, *models.EmployeeRole) (*models.EmployeeRole, *types.Error)

	UpdateStatus(*gin.Context, string, string) (*models.EmployeeRole, *types.Error)
}

type PermissionRepository interface {
	FindAll(*gin.Context, models.FindAllEmployeeRolePermissionParams) ([]*models.EmployeeRolePermission, *types.Error)
	Find(*gin.Context, string) (*models.EmployeeRolePermission, *types.Error)
	Create(*gin.Context, *models.CreateUpdateEmployeeRolePermission) (*models.EmployeeRolePermission, *types.Error)
	DeleteByEmployeeRoleID(*gin.Context, string) *types.Error

	CreateBunch(*gin.Context, string, models.FindAllEmployeeRolePermissionParams) *types.Error
}
