package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type CompanyRepository struct {
	repository       data.GenericStorage
	statusRepository data.GenericStorage
}

func NewCompanyRepository(repository data.GenericStorage, statusRepository data.GenericStorage) CompanyRepository {
	return CompanyRepository{repository: repository, statusRepository: statusRepository}
}

// A function to get all Data that matches the filter provided
func (s CompanyRepository) FindAll(ctx *gin.Context, params models.FindAllCompanyParams) ([]*models.Company, *types.Error) {
	result := []*models.Company{}
	bulks := []*models.CompanyBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where += fmt.Sprintf(` AND companies.%s`, params.FindAllParams.StatusID)
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT
    companies.id, companies.name, companies.code, companies.logo_img_url,
    companies.status_id,
    status.name AS status_name
  FROM companies
  JOIN status ON companies.status_id = status.id
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":  params.FindAllParams.Size,
		"offset": ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".CompanyStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		for _, v := range bulks {
			result = append(result, &models.Company{
				ID:         v.ID,
				Name:       v.Name,
				Code:       v.Code,
				LogoImgURL: v.LogoImgURL,
				StatusID:   v.StatusID,
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
func (s CompanyRepository) Find(ctx *gin.Context, id string) (*models.Company, *types.Error) {
	result := models.Company{}
	bulks := []*models.CompanyBulk{}
	var err error

	query := `
  SELECT
    companies.id, companies.name, companies.code, companies.logo_img_url,
    companies.status_id,
    status.name AS status_name
  FROM companies
  JOIN status ON companies.status_id = status.id
  WHERE companies.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".CompanyStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.Company{
			ID:         v.ID,
			Name:       v.Name,
			Code:       v.Code,
			LogoImgURL: v.LogoImgURL,
			StatusID:   v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}
	} else {
		return nil, &types.Error{
			Path:       ".CompanyStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

// Inserts a new row of data
func (s CompanyRepository) Create(ctx *gin.Context, obj *models.Company) (*models.Company, *types.Error) {
	data := models.Company{}
	_, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".CompanyStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".CompanyStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

// Updates a row of data specified by the given ID inside the obj struct
func (s CompanyRepository) Update(ctx *gin.Context, obj *models.Company) (*models.Company, *types.Error) {
	data := models.Company{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".CompanyStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".CompanyStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s CompanyRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.Company, *types.Error) {
	data := models.Company{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".CompanyStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".CompanyStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}
