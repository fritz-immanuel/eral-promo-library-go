package userpermission

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(*gin.Context, models.FindAllUserPermissionParams) ([]*models.UserPermission, *types.Error)
	Find(*gin.Context, string) (*models.UserPermission, *types.Error)
	Create(*gin.Context, models.CreateUpdateUserPermission) (*models.UserPermission, *types.Error)
	DeleteByUserID(*gin.Context, string) *types.Error

	CreateBunch(*gin.Context, string, models.FindAllUserPermissionParams) *types.Error
}
