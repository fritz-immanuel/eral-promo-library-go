package employeerole

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(context *gin.Context, params models.FindAllEmployeeRoleParams) ([]*models.EmployeeRole, *types.Error)
	Find(context *gin.Context, employeeroleID string) (*models.EmployeeRole, *types.Error)
	Count(context *gin.Context, params models.FindAllEmployeeRoleParams) (int, *types.Error)
	Create(context *gin.Context, data models.EmployeeRole) (*models.EmployeeRole, *types.Error)
	Update(context *gin.Context, employeeroleID string, data models.EmployeeRole) (*models.EmployeeRole, *types.Error)

	UpdateStatus(context *gin.Context, employeeroleID string, newStatusID string) (*models.EmployeeRole, *types.Error)
}
