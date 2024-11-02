package usecase

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/business"
	"github.com/google/uuid"

	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type BusinessUsecase struct {
	businessRepo   business.Repository
	contextTimeout time.Duration
	db             *sqlx.DB
}

func NewBusinessUsecase(db *sqlx.DB, businessRepo business.Repository) business.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &BusinessUsecase{
		businessRepo:   businessRepo,
		contextTimeout: timeoutContext,
		db:             db,
	}
}

func (u *BusinessUsecase) FindAll(ctx *gin.Context, params models.FindAllBusinessParams) ([]*models.Business, *types.Error) {
	result, err := u.businessRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".BusinessUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BusinessUsecase) Find(ctx *gin.Context, id string) (*models.Business, *types.Error) {
	result, err := u.businessRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".BusinessUsecase->Find()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BusinessUsecase) Count(ctx *gin.Context, params models.FindAllBusinessParams) (int, *types.Error) {
	result, err := u.businessRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".BusinessUsecase->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}

func (u *BusinessUsecase) Create(ctx *gin.Context, obj models.Business) (*models.Business, *types.Error) {
	data := models.Business{
		ID:         uuid.New().String(),
		Name:       obj.Name,
		Code:       obj.Code,
		LogoImgURL: obj.LogoImgURL,
		CompanyID:  obj.CompanyID,
		StatusID:   models.DEFAULT_STATUS_CODE,
	}

	result, err := u.businessRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".BusinessUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BusinessUsecase) Update(ctx *gin.Context, id string, obj models.Business) (*models.Business, *types.Error) {
	data, err := u.businessRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".BusinessUsecase->Update()" + err.Path
		return nil, err
	}

	data.Name = obj.Name
	data.Code = obj.Code
	data.LogoImgURL = obj.LogoImgURL
	data.CompanyID = obj.CompanyID

	result, err := u.businessRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".BusinessUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *BusinessUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.Business, *types.Error) {
	if newStatusID != models.STATUS_ACTIVE && newStatusID != models.STATUS_INACTIVE {
		return nil, &types.Error{
			Path:       ".BusinessUsecase->UpdateStatus()",
			Message:    "StatusID is not valid",
			Error:      fmt.Errorf("StatusID is not valid"),
			StatusCode: http.StatusBadRequest,
		}
	}

	result, err := u.businessRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".BusinessUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, nil
}
