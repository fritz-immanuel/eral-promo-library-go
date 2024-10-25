package businessconfig

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type Repository interface {
	FindAll(*gin.Context, models.FindAllBusinessConfigParams) ([]*models.BusinessConfig, *types.Error)
	Find(*gin.Context, int) (*models.BusinessConfig, *types.Error)
	Create(*gin.Context, *models.BusinessConfig) (*models.BusinessConfig, *types.Error)
	Update(*gin.Context, *models.BusinessConfig) (*models.BusinessConfig, *types.Error)
}
