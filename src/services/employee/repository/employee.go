package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type EmployeeRepository struct {
	repository       data.GenericStorage
	statusRepository data.GenericStorage
}

func NewEmployeeRepository(repository data.GenericStorage, statusRepository data.GenericStorage) EmployeeRepository {
	return EmployeeRepository{repository: repository, statusRepository: statusRepository}
}

func (s EmployeeRepository) FindAll(ctx *gin.Context, params models.FindAllEmployeeParams) ([]*models.Employee, *types.Error) {
	result := []*models.Employee{}
	bulks := []*models.EmployeeBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.StatusID)
	}

	if params.Email != "" {
		where += ` AND employees.email = :email`
	}

	if params.Username != "" {
		where += ` AND employees.username = :username`
	}

	if params.Password != "" {
		where += ` AND employees.password = :password`
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT
    employees.id, employees.name, employees.email, employees.password, employees.username,
    employees.status_id, status.name AS status_name
  FROM employees
  JOIN status ON status.id = employees.status_id
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":     params.FindAllParams.Size,
		"offset":    ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id": params.FindAllParams.StatusID,
		"email":     params.Email,
		"username":  params.Username,
		"password":  params.Password,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.Employee{
			ID:       v.ID,
			Name:     v.Name,
			Email:    v.Email,
			Username: v.Username,
			Password: v.Password,
			StatusID: v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}

		result = append(result, obj)
	}

	return result, nil
}

func (s EmployeeRepository) Find(ctx *gin.Context, id string) (*models.Employee, *types.Error) {
	var err error

	result := models.Employee{}
	bulks := []*models.EmployeeBulk{}

	query := `
  SELECT
    employees.id, employees.name, employees.email, employees.password, employees.username,
    employees.status_id, status.name AS status_name
  FROM employees
  JOIN status on status.id = employees.status_id
  WHERE employees.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.Employee{
			ID:       v.ID,
			Name:     v.Name,
			Email:    v.Email,
			Username: v.Username,
			Password: v.Password,
			StatusID: v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}
	} else {
		return nil, &types.Error{
			Path:       ".EmployeeStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s EmployeeRepository) Create(ctx *gin.Context, obj *models.Employee) (*models.Employee, *types.Error) {
	result := models.Employee{}
	_, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &result, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s EmployeeRepository) Update(ctx *gin.Context, obj *models.Employee) (*models.Employee, *types.Error) {
	result := models.Employee{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &result, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s EmployeeRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.Employee, *types.Error) {
	data := models.Employee{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}
