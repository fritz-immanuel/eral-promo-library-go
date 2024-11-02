package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type EmployeeRoleRepository struct {
	repository       data.GenericStorage
	statusRepository data.GenericStorage
}

func NewEmployeeRoleRepository(repository data.GenericStorage, statusRepository data.GenericStorage) EmployeeRoleRepository {
	return EmployeeRoleRepository{repository: repository, statusRepository: statusRepository}
}

func (s EmployeeRoleRepository) FindAll(ctx *gin.Context, params models.FindAllEmployeeRoleParams) ([]*models.EmployeeRole, *types.Error) {
	result := []*models.EmployeeRole{}
	bulks := []*models.EmployeeRoleBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.StatusID)
	}

	if params.IsSupervisor != 0 {
		where += ` AND employee_roles.is_supervisor = 1`
	}

	if params.IsNotSupervisor != 0 {
		where += ` AND employee_roles.is_supervisor = 0`
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT
    employee_roles.id, employee_roles.name, employee_roles.is_supervisor,
    employee_roles.status_id, status.name AS status_name
  FROM employee_roles
  JOIN status ON status.id = employee_roles.status_id
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":     params.FindAllParams.Size,
		"offset":    ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id": params.FindAllParams.StatusID,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeRoleStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.EmployeeRole{
			ID:           v.ID,
			Name:         v.Name,
			IsSupervisor: v.IsSupervisor,
			StatusID:     v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}

		result = append(result, obj)
	}

	return result, nil
}

func (s EmployeeRoleRepository) Find(ctx *gin.Context, id string) (*models.EmployeeRole, *types.Error) {
	var err error

	result := models.EmployeeRole{}
	bulks := []*models.EmployeeRoleBulk{}

	query := `
  SELECT
    employee_roles.id, employee_roles.name, employee_roles.is_supervisor,
    employee_roles.status_id, status.name AS status_name
  FROM employee_roles
  JOIN status ON status.id = employee_roles.status_id
  WHERE employee_roles.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeRoleStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.EmployeeRole{
			ID:           v.ID,
			Name:         v.Name,
			IsSupervisor: v.IsSupervisor,
			StatusID:     v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}
	} else {
		return nil, &types.Error{
			Path:       ".EmployeeRoleStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s EmployeeRoleRepository) Create(ctx *gin.Context, obj *models.EmployeeRole) (*models.EmployeeRole, *types.Error) {
	data := models.EmployeeRole{}
	_, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeRoleStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeRoleStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s EmployeeRoleRepository) Update(ctx *gin.Context, obj *models.EmployeeRole) (*models.EmployeeRole, *types.Error) {
	data := models.EmployeeRole{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeRoleStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeRoleStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s EmployeeRoleRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.EmployeeRole, *types.Error) {
	data := models.EmployeeRole{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeRoleStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeRoleStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}
