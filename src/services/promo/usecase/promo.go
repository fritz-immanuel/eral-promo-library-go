package usecase

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/promo"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/useraction"
	"github.com/google/uuid"

	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"

	useractionRepo "github.com/fritz-immanuel/eral-promo-library-go/src/services/useraction/repository"
)

type PromoUsecase struct {
	promoRepo         promo.Repository
	promodocumentRepo promo.DocumentRepository
	useractionRepo    useraction.Repository
	contextTimeout    time.Duration
	db                *sqlx.DB
}

func NewPromoUsecase(db *sqlx.DB, promoRepo promo.Repository, promodocumentRepo promo.DocumentRepository) promo.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	useractionRepo := useractionRepo.NewUserActionRepository(
		data.NewMySQLStorage(db, "user_actions", models.UserAction{}, data.MysqlConfig{}),
	)

	return &PromoUsecase{
		promoRepo:         promoRepo,
		promodocumentRepo: promodocumentRepo,
		useractionRepo:    useractionRepo,
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
		StartDate:        obj.StartDate,
		EndDate:          obj.EndDate,
		ImgURL:           obj.ImgURL,
		CompanyID:        obj.CompanyID,
		BusinessID:       obj.BusinessID,
		BrandID:          obj.BrandID,
		TotalPromoBudget: obj.TotalPromoBudget,
		PrincipleSupport: obj.PrincipleSupport,
		InternalSupport:  obj.InternalSupport,
		Description:      obj.Description,
		StatusID:         models.DEFAULT_STATUS_CODE,
	}

	result, err := u.promoRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".PromoUsecase->Create()" + err.Path
		return nil, err
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
	data.StartDate = obj.StartDate
	data.EndDate = obj.EndDate
	data.ImgURL = obj.ImgURL
	data.BrandID = obj.BrandID
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

// APPROVAL

func (u *PromoUsecase) ApprovePromo(ctx *gin.Context, id string) (*models.Promo, *types.Error) {
	result, err := u.promoRepo.ApprovePromo(ctx, id)
	if err != nil {
		err.Path = ".PromoUsecase->ApprovePromo()" + err.Path
		return nil, err
	}

	userAction := models.UserAction{}
	userAction.TableName = "promos"
	userAction.Action = "Approve Promo"
	userAction.RefID = id

	err = u.useractionRepo.CreateManual(ctx, &userAction)
	if err != nil {
		err.Path = ".PromoUsecase->ApprovePromo()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *PromoUsecase) RejectPromo(ctx *gin.Context, id string, rejectReason string) (*models.Promo, *types.Error) {
	result, err := u.promoRepo.RejectPromo(ctx, id, rejectReason)
	if err != nil {
		err.Path = ".PromoUsecase->RejectPromo()" + err.Path
		return nil, err
	}

	userAction := models.UserAction{}
	userAction.TableName = "promos"
	userAction.Action = "Reject Promo"
	userAction.RefID = id

	err = u.useractionRepo.CreateManual(ctx, &userAction)
	if err != nil {
		err.Path = ".PromoUsecase->ApprovePromo()" + err.Path
		return nil, err
	}

	return result, nil
}

// DOCUMENTS

func (u *PromoUsecase) FindDocument(ctx *gin.Context, id string) (*models.PromoDocument, *types.Error) {
	result, err := u.promodocumentRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".PromoUsecase->FindDocument()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *PromoUsecase) CreateDocument(ctx *gin.Context, obj models.PromoDocument) (*models.PromoDocument, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".PromoUsecase->CreateDocument()" + err.Path
		return nil, err
	}

	data := models.PromoDocument{
		ID:          uuid.New().String(),
		PromoID:     obj.PromoID,
		DocumentURL: obj.DocumentURL,
		StatusID:    models.DEFAULT_STATUS_CODE,
	}

	result, err := u.promodocumentRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".PromoUsecase->CreateDocument()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *PromoUsecase) UpdateDocument(ctx *gin.Context, id string, obj models.PromoDocument) (*models.PromoDocument, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".PromoUsecase->UpdateDocument()" + err.Path
		return nil, err
	}

	data, err := u.promodocumentRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".PromoUsecase->UpdateDocument()" + err.Path
		return nil, err
	}

	if data.PromoID != obj.PromoID {
		err = &types.Error{
			Path:       ".PromoUsecase->UpdateDocument()",
			Message:    "Data not found",
			Error:      fmt.Errorf("Promo ID does not match Document Promo ID"),
			StatusCode: http.StatusNotFound,
			Type:       "validation-error",
		}
		return nil, err
	}

	data.DocumentURL = obj.DocumentURL

	result, err := u.promodocumentRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".PromoUsecase->UpdateDocument()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *PromoUsecase) DeleteDocument(ctx *gin.Context, id string) *types.Error {
	_, err := u.promodocumentRepo.UpdateStatus(ctx, id, models.STATUS_INACTIVE)
	if err != nil {
		err.Path = ".PromoUsecase->DeleteDocument()" + err.Path
		return err
	}

	return nil
}

// HISTORY

func (u *PromoUsecase) FindUserActionHistory(ctx *gin.Context, id string, params models.FindAllActionHistory) ([]*models.UserAction, *types.Error) {
	filterFindAllParams := models.FindAllActionHistory{}
	filterFindAllParams = params
	filterFindAllParams.TableName = "promos"
	filterFindAllParams.UsingStatusTable = 1
	filterFindAllParams.RefID = id
	filterFindAllParams.FindAllParams.SortBy = "user_actions.created_at DESC"

	resultAction, err := u.useractionRepo.FindAll(ctx, filterFindAllParams)
	if err != nil {
		err.Path = ".PromoUsecase->FindUserActionHistory()" + err.Path
		return nil, err
	}

	return resultAction, err
}
