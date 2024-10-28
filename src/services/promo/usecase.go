package promo

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(context *gin.Context, params models.FindAllPromoParams) ([]*models.Promo, *types.Error)
	Find(context *gin.Context, id string) (*models.Promo, *types.Error)
	Count(context *gin.Context, params models.FindAllPromoParams) (int, *types.Error)
	Create(context *gin.Context, newData models.Promo) (*models.Promo, *types.Error)
	Update(context *gin.Context, id string, updatedData models.Promo) (*models.Promo, *types.Error)

	FindStatus(context *gin.Context) ([]*models.Status, *types.Error)
	UpdateStatus(*gin.Context, string, string) (*models.Promo, *types.Error)

	// APPROVAL
	ApprovePromo(*gin.Context, string) (*models.Promo, *types.Error)
	RejectPromo(*gin.Context, string, string) (*models.Promo, *types.Error)

	// DOCUMENT
	FindDocument(context *gin.Context, id string) (*models.PromoDocument, *types.Error)
	CreateDocument(*gin.Context, models.PromoDocument) (*models.PromoDocument, *types.Error)
	UpdateDocument(*gin.Context, string, models.PromoDocument) (*models.PromoDocument, *types.Error)
	DeleteDocument(*gin.Context, string) *types.Error

	// HISTORY
	FindUserActionHistory(*gin.Context, string, models.FindAllActionHistory) ([]*models.UserAction, *types.Error)
}
