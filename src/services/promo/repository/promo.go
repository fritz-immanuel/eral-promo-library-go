package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type PromoRepository struct {
	repository       data.GenericStorage
	statusRepository data.GenericStorage
}

func NewPromoRepository(repository data.GenericStorage, statusRepository data.GenericStorage) PromoRepository {
	return PromoRepository{repository: repository, statusRepository: statusRepository}
}

// A function to get all Data that matches the filter provided
func (s PromoRepository) FindAll(ctx *gin.Context, params models.FindAllPromoParams) ([]*models.Promo, *types.Error) {
	result := []*models.Promo{}
	bulks := []*models.PromoBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where += fmt.Sprintf(` AND promos.%s`, params.FindAllParams.StatusID)
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT
    promos.id, promos.name, promos.code, promos.promo_type_id, promos.img_url, promos.start_date, promos.end_date,
    promos.company_id, promos.business_id, promos.total_promo_budget, promos.principle_support, promos.internal_support,
    promos.description,
    promos.status_id, promo_status.name AS status_name
  FROM promos
  JOIN promo_status ON promos.status_id = promo_status.id
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":  params.FindAllParams.Size,
		"offset": ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		for _, v := range bulks {
			result = append(result, &models.Promo{
				ID:               v.ID,
				Name:             v.Name,
				Code:             v.Code,
				PromoTypeID:      v.PromoTypeID,
				StartDate:        v.StartDate,
				EndDate:          v.EndDate,
				ImgURL:           v.ImgURL,
				CompanyID:        v.CompanyID,
				BusinessID:       v.BusinessID,
				TotalPromoBudget: v.TotalPromoBudget,
				PrincipleSupport: v.PrincipleSupport,
				InternalSupport:  v.InternalSupport,
				Description:      v.Description,
				StatusID:         v.StatusID,
				Status: models.Status{
					ID:   v.StatusID,
					Name: v.StatusName,
				},
			})
		}
	}

	return result, nil
}

// A function to get a row of data specified by the given ID
func (s PromoRepository) Find(ctx *gin.Context, id string) (*models.Promo, *types.Error) {
	result := models.Promo{}
	bulks := []*models.PromoBulk{}
	var err error

	query := `
  SELECT
    promos.id, promos.name, promos.code, promos.promo_type_id, promos.img_url, promos.start_date, promos.end_date,
    promos.company_id, promos.business_id, promos.total_promo_budget, promos.principle_support, promos.internal_support,
    promos.description,
    promos.status_id, promo_status.name AS status_name
  FROM promos
  JOIN promo_status ON promos.status_id = promo_status.id
  WHERE promos.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.Promo{
			ID:               v.ID,
			Name:             v.Name,
			Code:             v.Code,
			PromoTypeID:      v.PromoTypeID,
			StartDate:        v.StartDate,
			EndDate:          v.EndDate,
			ImgURL:           v.ImgURL,
			CompanyID:        v.CompanyID,
			BusinessID:       v.BusinessID,
			TotalPromoBudget: v.TotalPromoBudget,
			PrincipleSupport: v.PrincipleSupport,
			InternalSupport:  v.InternalSupport,
			Description:      v.Description,
			StatusID:         v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}
	} else {
		return nil, &types.Error{
			Path:       ".PromoStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

// Inserts a new row of data
func (s PromoRepository) Create(ctx *gin.Context, obj *models.Promo) (*models.Promo, *types.Error) {
	result := models.Promo{}
	_, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &result, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

// Updates a row of data specified by the given ID inside the obj struct
func (s PromoRepository) Update(ctx *gin.Context, obj *models.Promo) (*models.Promo, *types.Error) {
	result := models.Promo{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &result, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s PromoRepository) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	status := []*models.Status{}

	err := s.statusRepository.Where(ctx, &status, "1=1", map[string]interface{}{})

	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoStorage->FindStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return status, nil
}

func (s PromoRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.Promo, *types.Error) {
	data := models.Promo{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}
