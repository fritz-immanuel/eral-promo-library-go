package usecase

import (
	"time"

	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/company"
	"github.com/google/uuid"

	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type CompanyUsecase struct {
	companyRepo    company.Repository
	contextTimeout time.Duration
	db             *sqlx.DB
}

func NewCompanyUsecase(db *sqlx.DB, companyRepo company.Repository) company.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &CompanyUsecase{
		companyRepo:    companyRepo,
		contextTimeout: timeoutContext,
		db:             db,
	}
}

func (u *CompanyUsecase) FindAll(ctx *gin.Context, params models.FindAllCompanyParams) ([]*models.Company, *types.Error) {
	result, err := u.companyRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".CompanyUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *CompanyUsecase) Find(ctx *gin.Context, id string) (*models.Company, *types.Error) {
	result, err := u.companyRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".CompanyUsecase->Find()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *CompanyUsecase) Count(ctx *gin.Context, params models.FindAllCompanyParams) (int, *types.Error) {
	result, err := u.companyRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".CompanyUsecase->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}

func (u *CompanyUsecase) Create(ctx *gin.Context, obj models.Company) (*models.Company, *types.Error) {
	data := models.Company{
		ID:         uuid.New().String(),
		Name:       obj.Name,
		Code:       obj.Code,
		LogoImgURL: obj.LogoImgURL,
		StatusID:   models.DEFAULT_STATUS_CODE,
	}

	result, err := u.companyRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".CompanyUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *CompanyUsecase) Update(ctx *gin.Context, id string, obj models.Company) (*models.Company, *types.Error) {
	data, err := u.companyRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".CompanyUsecase->Update()" + err.Path
		return nil, err
	}

	data.Name = obj.Name
	data.Code = obj.Code
	data.LogoImgURL = obj.LogoImgURL

	result, err := u.companyRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".CompanyUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *CompanyUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.Company, *types.Error) {
	result, err := u.companyRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".CompanyUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, nil
}
