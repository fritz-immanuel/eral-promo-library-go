package company

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(context *gin.Context, params models.FindAllCompanyParams) ([]*models.Company, *types.Error)
	Find(context *gin.Context, id string) (*models.Company, *types.Error)
	Count(context *gin.Context, params models.FindAllCompanyParams) (int, *types.Error)
	Create(context *gin.Context, newData models.Company) (*models.Company, *types.Error)
	Update(context *gin.Context, id string, updatedData models.Company) (*models.Company, *types.Error)

	UpdateStatus(*gin.Context, string, string) (*models.Company, *types.Error)
}
