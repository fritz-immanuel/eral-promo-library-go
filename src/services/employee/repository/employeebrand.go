package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type EmployeeBrandRepository struct {
	repository data.GenericStorage
}

func NewEmployeeBrandRepository(repository data.GenericStorage) EmployeeBrandRepository {
	return EmployeeBrandRepository{repository: repository}
}

func (s EmployeeBrandRepository) FindAll(ctx *gin.Context, params models.FindAllEmployeeBrandParams) ([]*models.EmployeeBrand, *types.Error) {
	result := []*models.EmployeeBrand{}
	bulks := []*models.EmployeeBrandBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.EmployeeID != "" {
		where += ` AND employee_brands.employee_id = :employee_id`
	}

	if params.BrandID != "" {
		where += ` AND employee_brands.brand_id = :brand_id`
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT
    employee_brands.id, employee_brands.employee_id, employee_brands.brand_id,
		brands.name brand_name, brands.code brand_code
  FROM employee_brands
	JOIN brands ON brands.id = employee_brands.brand_id
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":       params.FindAllParams.Size,
		"offset":      ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id":   params.FindAllParams.StatusID,
		"employee_id": params.EmployeeID,
		"brand_id":    params.BrandID,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeBrandStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.EmployeeBrand{
			ID:         v.ID,
			EmployeeID: v.EmployeeID,
			BrandID:    v.BrandID,
			Brand: &models.StringIDNameCodeTemplate{
				ID:   v.BrandID,
				Name: v.BrandName,
				Code: v.BrandCode,
			},
		}

		result = append(result, obj)
	}

	return result, nil
}

func (s EmployeeBrandRepository) Find(ctx *gin.Context, id string) (*models.EmployeeBrand, *types.Error) {
	var err error

	result := models.EmployeeBrand{}
	bulks := []*models.EmployeeBrandBulk{}

	query := `
  SELECT
    employee_brands.id, employee_brands.employee_id, employee_brands.brand_id,
		brands.name brand_name, brands.code brand_code
  FROM employee_brands
	JOIN brands ON brands.id = employee_brands.brand_id
  WHERE employee_brands.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeBrandStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.EmployeeBrand{
			ID:         v.ID,
			EmployeeID: v.EmployeeID,
			BrandID:    v.BrandID,
			Brand: &models.StringIDNameCodeTemplate{
				ID:   v.BrandID,
				Name: v.BrandName,
				Code: v.BrandCode,
			},
		}
	} else {
		return nil, &types.Error{
			Path:       ".EmployeeBrandStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s EmployeeBrandRepository) Create(ctx *gin.Context, obj *models.EmployeeBrand) (*models.EmployeeBrand, *types.Error) {
	result := models.EmployeeBrand{}
	_, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeBrandStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &result, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeBrandStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s EmployeeBrandRepository) Update(ctx *gin.Context, obj *models.EmployeeBrand) (*models.EmployeeBrand, *types.Error) {
	result := models.EmployeeBrand{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeBrandStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &result, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeBrandStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s EmployeeBrandRepository) DeleteByEmployeeID(ctx *gin.Context, id string) *types.Error {
	args := make(map[string]interface{})
	err := s.repository.ExecQuery(ctx, fmt.Sprintf("DELETE FROM employee_brands WHERE employee_id = '%s'", id), args)
	if err != nil {
		return &types.Error{
			Path:       ".EmployeeBrandStorage->DeleteByEmployeeID()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return nil
}
