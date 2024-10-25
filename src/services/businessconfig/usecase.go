package businessconfig

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type Usecase interface {
	FindAll(*gin.Context, models.FindAllBusinessConfigParams) ([]*models.BusinessConfig, *types.Error)
	Find(*gin.Context, int) (*models.BusinessConfig, *types.Error)
	Count(*gin.Context, models.FindAllBusinessConfigParams) (int, *types.Error)
	Create(*gin.Context, models.BusinessConfig) (*models.BusinessConfig, *types.Error)
	Update(*gin.Context, int, models.BusinessConfig) (*models.BusinessConfig, *types.Error)
}
