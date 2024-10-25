package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type BusinessConfigRepository struct {
	repository data.GenericStorage
}

func NewBusinessConfigRepository(repository data.GenericStorage) BusinessConfigRepository {
	return BusinessConfigRepository{repository: repository}
}

func (s BusinessConfigRepository) FindAll(ctx *gin.Context, params models.FindAllBusinessConfigParams) ([]*models.BusinessConfig, *types.Error) {
	result := []*models.BusinessConfig{}
	bulks := []*models.BusinessConfigBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.StatusID)
	}

	if params.BusinessID != 0 {
		where += fmt.Sprintf(` AND business_configs.business_id = %d`, params.BusinessID)
	}

	if params.SubURLName != "" {
		where += ` AND business_configs.sub_url_name = :sub_url_name`
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT
    business_configs.id, business_configs.business_id,
		business_configs.sub_url_name, business_configs.config
  FROM business_configs
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":        params.FindAllParams.Size,
		"offset":       ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id":    params.FindAllParams.StatusID,
		"sub_url_name": params.SubURLName,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".BusinessConfigStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.BusinessConfig{
			ID:         v.ID,
			BusinessID: v.BusinessID,
			SubURLName: v.SubURLName,
			Config:     v.Config,
		}

		result = append(result, obj)
	}

	return result, nil
}

func (s BusinessConfigRepository) Find(ctx *gin.Context, id int) (*models.BusinessConfig, *types.Error) {
	var err error

	result := models.BusinessConfig{}
	bulks := []*models.BusinessConfigBulk{}

	query := fmt.Sprintf(`
  SELECT
    business_configs.id, business_configs.business_id,
		business_configs.sub_url_name, business_configs.config
  FROM business_configs
  WHERE business_configs.id = %d`, id)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{})
	if err != nil {
		return nil, &types.Error{
			Path:       ".BusinessConfigStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.BusinessConfig{
			ID:         v.ID,
			BusinessID: v.BusinessID,
			SubURLName: v.SubURLName,
			Config:     v.Config,
		}
	} else {
		return nil, &types.Error{
			Path:       ".BusinessConfigStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s BusinessConfigRepository) Create(ctx *gin.Context, obj *models.BusinessConfig) (*models.BusinessConfig, *types.Error) {
	data := models.BusinessConfig{}
	result, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BusinessConfigStorage->Create()",
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
			Path:       ".BusinessConfigStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s BusinessConfigRepository) Update(ctx *gin.Context, obj *models.BusinessConfig) (*models.BusinessConfig, *types.Error) {
	data := models.BusinessConfig{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BusinessConfigStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BusinessConfigStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}
