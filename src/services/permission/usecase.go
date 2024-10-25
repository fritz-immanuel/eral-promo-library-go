package permission

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(*gin.Context, models.FindAllPermissionParams) ([]*models.Permission, *types.Error)
	Find(*gin.Context, int) (*models.Permission, *types.Error)
	Count(*gin.Context, models.FindAllPermissionParams) (int, *types.Error)
}
