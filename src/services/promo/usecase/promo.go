package usecase

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/promo"
	"github.com/google/uuid"

	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type PromoUsecase struct {
	promoRepo         promo.Repository
	promodocumentRepo promo.DocumentRepository
	contextTimeout    time.Duration
	db                *sqlx.DB
}

func NewPromoUsecase(db *sqlx.DB, promoRepo promo.Repository, promodocumentRepo promo.DocumentRepository) promo.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &PromoUsecase{
		promoRepo:         promoRepo,
		promodocumentRepo: promodocumentRepo,
		contextTimeout:    timeoutContext,
		db:                db,
	}
}

func (u *PromoUsecase) FindAll(ctx *gin.Context, params models.FindAllPromoParams) ([]*models.Promo, *types.Error) {
	result, err := u.promoRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".PromoUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *PromoUsecase) Find(ctx *gin.Context, id string) (*models.Promo, *types.Error) {
	result, err := u.promoRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".PromoUsecase->Find()" + err.Path
		return nil, err
	}

	var documentParams models.FindAllPromoDocumentParams
	documentParams.FindAllParams.StatusID = `status_id = 1`
	documentParams.PromoID = id
	result.PromoDocuments, err = u.promodocumentRepo.FindAll(ctx, documentParams)
	if err != nil {
		err.Path = ".PromoUsecase->Find()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *PromoUsecase) Count(ctx *gin.Context, params models.FindAllPromoParams) (int, *types.Error) {
	result, err := u.promoRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".PromoUsecase->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}

func (u *PromoUsecase) Create(ctx *gin.Context, obj models.Promo) (*models.Promo, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".PromoUsecase->Create()" + err.Path
		return nil, err
	}

	if obj.PrincipleSupport > 100.0 || obj.InternalSupport > 100.0 {
		err = &types.Error{
			Path:       ".PromoUsecase->Create()",
			Message:    "Principle Support and Internal Support must be less than 100",
			Error:      fmt.Errorf("PrincipleSupport and InternalSupport must be less than 100"),
			StatusCode: http.StatusBadRequest,
			Type:       "validation-error",
		}
		return nil, err
	}

	if obj.PrincipleSupport+obj.InternalSupport > 100.0 {
		err = &types.Error{
			Path:       ".PromoUsecase->Create()",
			Message:    "Principle Support and Internal Support TOTAL must be less than 100",
			Error:      fmt.Errorf("PrincipleSupport and InternalSupport TOTAL must be less than 100"),
			StatusCode: http.StatusBadRequest,
			Type:       "validation-error",
		}
		return nil, err
	}

	data := models.Promo{
		ID:               uuid.New().String(),
		Name:             obj.Name,
		Code:             obj.Code,
		PromoTypeID:      obj.PromoTypeID,
		StartDate:        obj.StartDate,
		EndDate:          obj.EndDate,
		ImgURL:           obj.ImgURL,
		CompanyID:        obj.CompanyID,
		BusinessID:       obj.BusinessID,
		TotalPromoBudget: obj.TotalPromoBudget,
		PrincipleSupport: obj.PrincipleSupport,
		InternalSupport:  obj.InternalSupport,
		Description:      obj.Description,
		StatusID:         obj.StatusID,
	}

	result, err := u.promoRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".PromoUsecase->Create()" + err.Path
		return nil, err
	}

	if len(obj.PromoDocuments) > 0 {
		for _, v := range obj.PromoDocuments {
			v.ID = uuid.New().String()
			v.PromoID = data.ID
			_, err := u.promodocumentRepo.Create(ctx, v)
			if err != nil {
				err.Path = ".PromoUsecase->Create()" + err.Path
				return nil, err
			}
		}
	}

	return result, nil
}

func (u *PromoUsecase) Update(ctx *gin.Context, id string, obj models.Promo) (*models.Promo, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".PromoUsecase->Update()" + err.Path
		return nil, err
	}

	if obj.PrincipleSupport > 100.0 || obj.InternalSupport > 100.0 {
		err = &types.Error{
			Path:       ".PromoUsecase->Update()",
			Message:    "Principle Support and Internal Support must be less than 100",
			Error:      fmt.Errorf("PrincipleSupport and InternalSupport must be less than 100"),
			StatusCode: http.StatusBadRequest,
			Type:       "validation-error",
		}
		return nil, err
	}

	if obj.PrincipleSupport+obj.InternalSupport > 100.0 {
		err = &types.Error{
			Path:       ".PromoUsecase->Update()",
			Message:    "Principle Support and Internal Support TOTAL must be less than 100",
			Error:      fmt.Errorf("PrincipleSupport and InternalSupport TOTAL must be less than 100"),
			StatusCode: http.StatusBadRequest,
			Type:       "validation-error",
		}
		return nil, err
	}

	data, err := u.promoRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".PromoUsecase->Update()" + err.Path
		return nil, err
	}

	data.Name = obj.Name
	data.Code = obj.Code
	data.PromoTypeID = obj.PromoTypeID
	data.StartDate = obj.StartDate
	data.EndDate = obj.EndDate
	data.ImgURL = obj.ImgURL
	data.CompanyID = obj.CompanyID
	data.BusinessID = obj.BusinessID
	data.TotalPromoBudget = obj.TotalPromoBudget
	data.PrincipleSupport = obj.PrincipleSupport
	data.InternalSupport = obj.InternalSupport
	data.Description = obj.Description

	result, err := u.promoRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".PromoUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *PromoUsecase) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	result, err := u.promoRepo.FindStatus(ctx)
	if err != nil {
		err.Path = ".PromoUsecase->FindStatus()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *PromoUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.Promo, *types.Error) {
	result, err := u.promoRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".PromoUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, nil
}
