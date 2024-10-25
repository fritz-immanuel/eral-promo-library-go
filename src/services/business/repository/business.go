package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type BusinessRepository struct {
	repository       data.GenericStorage
	statusRepository data.GenericStorage
}

func NewBusinessRepository(repository data.GenericStorage, statusRepository data.GenericStorage) BusinessRepository {
	return BusinessRepository{repository: repository, statusRepository: statusRepository}
}

// A function to get all Data that matches the filter provided
func (s BusinessRepository) FindAll(ctx *gin.Context, params models.FindAllBusinessParams) ([]*models.Business, *types.Error) {
	result := []*models.Business{}
	bulks := []*models.BusinessBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where += fmt.Sprintf(` AND business.%s`, params.FindAllParams.StatusID)
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT
    business.id, business.name, business.code, business.logo_img,
    business.status_id,
    status.name AS status_name
  FROM business
  JOIN status ON business.status_id = status.id
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":  params.FindAllParams.Size,
		"offset": ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".BusinessStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		for _, v := range bulks {
			result = append(result, &models.Business{
				ID:       v.ID,
				Name:     v.Name,
				Code:     v.Code,
				LogoImg:  v.LogoImg,
				StatusID: v.StatusID,
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
func (s BusinessRepository) Find(ctx *gin.Context, id string) (*models.Business, *types.Error) {
	result := models.Business{}
	bulks := []*models.BusinessBulk{}
	var err error

	query := `
  SELECT
    business.id, business.name, business.code, business.logo_img,
    business.status_id,
    status.name AS status_name
  FROM business
  JOIN status ON business.status_id = status.id
  WHERE business.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".BusinessStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.Business{
			ID:       v.ID,
			Name:     v.Name,
			Code:     v.Code,
			LogoImg:  v.LogoImg,
			StatusID: v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}
	} else {
		return nil, &types.Error{
			Path:       ".BusinessStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

// Inserts a new row of data
func (s BusinessRepository) Create(ctx *gin.Context, obj *models.Business) (*models.Business, *types.Error) {
	data := models.Business{}
	result, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BusinessStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	lastID, _ := (*result).LastInsertId()
	err = s.repository.FindByID(ctx, &data, lastID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BusinessStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

// Updates a row of data specified by the given ID inside the obj struct
func (s BusinessRepository) Update(ctx *gin.Context, obj *models.Business) (*models.Business, *types.Error) {
	data := models.Business{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BusinessStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BusinessStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s BusinessRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.Business, *types.Error) {
	data := models.Business{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BusinessStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BusinessStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}
