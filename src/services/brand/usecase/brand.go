package usecase

import (
	"time"

	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/brand"
	"github.com/google/uuid"

	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type BrandUsecase struct {
	brandRepo      brand.Repository
	contextTimeout time.Duration
	db             *sqlx.DB
}

func NewBrandUsecase(db *sqlx.DB, brandRepo brand.Repository) brand.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &BrandUsecase{
		brandRepo:      brandRepo,
		contextTimeout: timeoutContext,
		db:             db,
	}
}

func (u *BrandUsecase) FindAll(ctx *gin.Context, params models.FindAllBrandParams) ([]*models.Brand, *types.Error) {
	result, err := u.brandRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".BrandUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BrandUsecase) Find(ctx *gin.Context, id string) (*models.Brand, *types.Error) {
	result, err := u.brandRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".BrandUsecase->Find()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BrandUsecase) Count(ctx *gin.Context, params models.FindAllBrandParams) (int, *types.Error) {
	result, err := u.brandRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".BrandUsecase->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}

func (u *BrandUsecase) Create(ctx *gin.Context, obj models.Brand) (*models.Brand, *types.Error) {
	data := models.Brand{
		ID:         uuid.New().String(),
		Name:       obj.Name,
		Code:       obj.Code,
		LogoImgURL: obj.LogoImgURL,
		BusinessID: obj.BusinessID,
		StatusID:   models.DEFAULT_STATUS_CODE,
	}

	result, err := u.brandRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".BrandUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BrandUsecase) Update(ctx *gin.Context, id string, obj models.Brand) (*models.Brand, *types.Error) {
	data, err := u.brandRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".BrandUsecase->Update()" + err.Path
		return nil, err
	}

	data.Name = obj.Name
	data.Code = obj.Code
	data.LogoImgURL = obj.LogoImgURL
	data.BusinessID = obj.BusinessID

	result, err := u.brandRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".BrandUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *BrandUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.Brand, *types.Error) {
	result, err := u.brandRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".BrandUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, nil
}
