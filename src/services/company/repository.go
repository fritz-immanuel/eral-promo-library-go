package company

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

// Repository is the contract between Repository and usecase
type Repository interface {
	FindAll(*gin.Context, models.FindAllCompanyParams) ([]*models.Company, *types.Error)
	Find(*gin.Context, string) (*models.Company, *types.Error)
	Create(*gin.Context, *models.Company) (*models.Company, *types.Error)
	Update(*gin.Context, *models.Company) (*models.Company, *types.Error)
	UpdateStatus(*gin.Context, string, string) (*models.Company, *types.Error)
}
