package user

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type Repository interface {
	FindAll(*gin.Context, models.FindAllUserParams) ([]*models.User, *types.Error)
	Find(*gin.Context, string) (*models.User, *types.Error)
	Create(*gin.Context, *models.User) (*models.User, *types.Error)
	Update(*gin.Context, *models.User) (*models.User, *types.Error)
	UpdateStatus(*gin.Context, string, string) (*models.User, *types.Error)
}
