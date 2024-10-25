package employee

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(context *gin.Context, params models.FindAllEmployeeParams) ([]*models.Employee, *types.Error)
	Find(context *gin.Context, employeeID string) (*models.Employee, *types.Error)
	Count(context *gin.Context, params models.FindAllEmployeeParams) (int, *types.Error)
	Create(context *gin.Context, data models.Employee) (*models.Employee, *types.Error)
	Update(context *gin.Context, employeeID string, data models.Employee) (*models.Employee, *types.Error)
	UpdateStatus(context *gin.Context, employeeID string, newStatusID string) (*models.Employee, *types.Error)

	// PASSWORD
	UpdatePassword(*gin.Context, string, string) (*models.Employee, *types.Error)

	Login(context *gin.Context, creds models.EmployeeLogin) (*models.EmployeeLogin, *types.Error)
}
