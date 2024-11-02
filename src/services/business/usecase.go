package business

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(context *gin.Context, params models.FindAllBusinessParams) ([]*models.Business, *types.Error)
	Find(context *gin.Context, id string) (*models.Business, *types.Error)
	Count(context *gin.Context, params models.FindAllBusinessParams) (int, *types.Error)
	Create(context *gin.Context, newData models.Business) (*models.Business, *types.Error)
	Update(context *gin.Context, id string, updatedData models.Business) (*models.Business, *types.Error)

	UpdateStatus(*gin.Context, string, string) (*models.Business, *types.Error)
}
