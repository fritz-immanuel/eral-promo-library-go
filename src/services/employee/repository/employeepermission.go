package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type EmployeePermissionRepository struct {
	repository data.GenericStorage
}

func NewEmployeePermissionRepository(repository data.GenericStorage) EmployeePermissionRepository {
	return EmployeePermissionRepository{repository: repository}
}

func (s EmployeePermissionRepository) FindAll(ctx *gin.Context, params models.FindAllEmployeePermissionParams) ([]*models.EmployeePermission, *types.Error) {
	result := []*models.EmployeePermission{}
	bulks := []*models.EmployeePermissionBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.EmployeeID != "" {
		where += ` AND employee_permissions.employee_id = :employee_id`
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT
    employee_permissions.employee_id,
    employee_permissions.permission_id,
    permission.package AS permission_package,
    permission.module_name AS permission_module_name,
    permission.action_name AS permission_action_name,
    permission.http_method AS permission_http_method,
    permission.route AS permission_route
  FROM employee_permissions
  JOIN permission on permission.id = employee_permissions.permission_id
  WHERE %s
  `, where)

	// fmt.Println(query)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"employee_id": params.EmployeeID,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeePermissionStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.EmployeePermission{
			EmployeeID:   v.EmployeeID,
			PermissionID: v.PermissionID,
			Permission: models.Permission{
				ID:         v.PermissionID,
				Package:    v.PermissionPackage,
				ModuleName: v.PermissionModuleName,
				ActionName: v.PermissionActionName,
				HTTPMethod: v.PermissionHTTPMethod,
				Route:      v.PermissionRoute,
			},
		}

		result = append(result, obj)
	}

	return result, nil
}

func (s EmployeePermissionRepository) Find(ctx *gin.Context, id string) (*models.EmployeePermission, *types.Error) {
	result := &models.EmployeePermission{}
	bulks := []*models.EmployeePermissionBulk{}

	var err error

	query := `
  SELECT
    employee_permissions.employee_id,
    employee_permissions.permission_id,
    permission.package AS permission_package,
    permission.module_name AS permission_module_name,
    permission.action_name AS permission_action_name,
    permission.display_module_name AS permission_display_module_name,
    permission.display_action_name AS permission_display_action_name,
    permission.http_method AS permission_http_method,
    permission.route AS permission_route
  FROM employee_permissions
  JOIN permission on permission.id = employee_permissions.permission_id
  WHERE employee_permissions.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{})
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeePermissionStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = &models.EmployeePermission{
			EmployeeID:   v.EmployeeID,
			PermissionID: v.PermissionID,
			Permission: models.Permission{
				ID:                v.PermissionID,
				Package:           v.PermissionPackage,
				ModuleName:        v.PermissionModuleName,
				ActionName:        v.PermissionActionName,
				DisplayModuleName: v.PermissionDisplayModuleName,
				DisplayActionName: v.PermissionDisplayActionName,
				HTTPMethod:        v.PermissionHTTPMethod,
				Route:             v.PermissionRoute,
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

	return result, nil
}

func (s EmployeePermissionRepository) Create(ctx *gin.Context, obj *models.CreateUpdateEmployeePermission) (*models.EmployeePermission, *types.Error) {
	data := models.EmployeePermission{}
	_, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeePermissionStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	data.EmployeeID = obj.EmployeeID
	data.PermissionID = obj.PermissionID
	return &data, nil
}

func (s EmployeePermissionRepository) DeleteByEmployeeID(ctx *gin.Context, id string) *types.Error {
	args := make(map[string]interface{})
	err := s.repository.ExecQuery(ctx, fmt.Sprintf("DELETE FROM employee_permissions WHERE employee_id = %s", id), args)
	if err != nil {
		return &types.Error{
			Path:       ".EmployeePermissionStorage->DeleteByEmployeeID()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return nil
}

func (s EmployeePermissionRepository) CreateBunch(ctx *gin.Context, employeeID string, params models.FindAllEmployeePermissionParams) *types.Error {
	args := make(map[string]interface{})

	where := "TRUE"
	not := ""

	if params.Package != "" {
		where = fmt.Sprintf("%s AND package = '%s'", where, params.Package)
	}

	if params.PermissionIDString != "" {
		where = fmt.Sprintf("%s AND id IN (%s)", where, params.PermissionIDString)
	}

	if params.Not != 0 {
		not = "NOT"
	}

	query := fmt.Sprintf(`
  INSERT INTO employee_permissions (employee_id, permission_id, created_at, updated_at)
  SELECT "%s", id, UTC_TIMESTAMP + INTERVAL 7 hour, UTC_TIMESTAMP + INTERVAL 7 HOUR
  FROM (
    SELECT id FROM permission
    WHERE %s AND id %s IN (
      SELECT permission_id FROM employee_permissions
      WHERE employee_id = "%s"
    )
  ) permission`, employeeID, where, not, params.EmployeeID)

	err := s.repository.ExecQuery(ctx, query, args)
	if err != nil {
		return &types.Error{
			Path:       ".EmployeePermissionStorage->CreateBunch()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return nil
}
