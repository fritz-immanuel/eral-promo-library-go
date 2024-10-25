package useraction

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(*gin.Context, models.FindAllActionHistory) ([]*models.UserAction, *types.Error)
	CreateManual(*gin.Context, models.UserAction) *types.Error
}
