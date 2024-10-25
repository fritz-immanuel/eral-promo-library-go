package usecase

import (
	"time"

	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/businessconfig"

	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type BusinessConfigUsecase struct {
	businessconfigRepo businessconfig.Repository
	contextTimeout     time.Duration
	db                 *sqlx.DB
}

func NewBusinessConfigUsecase(db *sqlx.DB, businessconfigRepo businessconfig.Repository) businessconfig.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &BusinessConfigUsecase{
		businessconfigRepo: businessconfigRepo,
		contextTimeout:     timeoutContext,
		db:                 db,
	}
}

func (u *BusinessConfigUsecase) FindAll(ctx *gin.Context, params models.FindAllBusinessConfigParams) ([]*models.BusinessConfig, *types.Error) {
	result, err := u.businessconfigRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".BusinessConfigUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BusinessConfigUsecase) Find(ctx *gin.Context, id int) (*models.BusinessConfig, *types.Error) {
	result, err := u.businessconfigRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".BusinessConfigUsecase->Find()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *BusinessConfigUsecase) Count(ctx *gin.Context, params models.FindAllBusinessConfigParams) (int, *types.Error) {
	result, err := u.businessconfigRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".BusinessConfigUsecase->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}

func (u *BusinessConfigUsecase) Create(ctx *gin.Context, obj models.BusinessConfig) (*models.BusinessConfig, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".BusinessConfigUsecase->Create()" + err.Path
		return nil, err
	}

	data := models.BusinessConfig{}
	data.ID = obj.ID
	data.BusinessID = obj.BusinessID
	data.SubURLName = obj.SubURLName
	data.Config = obj.Config

	result, err := u.businessconfigRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".BusinessConfigUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BusinessConfigUsecase) Update(ctx *gin.Context, id int, obj models.BusinessConfig) (*models.BusinessConfig, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".BusinessConfigUsecase->Update()" + err.Path
		return nil, err
	}

	data, err := u.businessconfigRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".BusinessConfigUsecase->Update()" + err.Path
		return nil, err
	}

	data.BusinessID = obj.BusinessID
	data.SubURLName = obj.SubURLName
	data.Config = obj.Config

	result, err := u.businessconfigRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".BusinessConfigUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}
