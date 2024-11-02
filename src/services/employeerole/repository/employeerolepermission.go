package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type EmployeeRolePermissionRepository struct {
	repository data.GenericStorage
}

func NewEmployeeRolePermissionRepository(repository data.GenericStorage) EmployeeRolePermissionRepository {
	return EmployeeRolePermissionRepository{repository: repository}
}

func (s EmployeeRolePermissionRepository) FindAll(ctx *gin.Context, params models.FindAllEmployeeRolePermissionParams) ([]*models.EmployeeRolePermission, *types.Error) {
	result := []*models.EmployeeRolePermission{}
	bulks := []*models.EmployeeRolePermissionBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.EmployeeRoleID != "" {
		where += ` AND employee_role_permissions.employee_role_id = :employee_role_id`
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT
    employee_role_permissions.employee_role_id,
    employee_role_permissions.permission_id,
    permissions.package AS permission_package,
    permissions.module_name AS permission_module_name,
    permissions.action_name AS permission_action_name,
    permissions.http_method AS permission_http_method,
    permissions.route AS permission_route
  FROM employee_role_permissions
  JOIN permissions on permissions.id = employee_role_permissions.permission_id
  WHERE %s
  `, where)

	// fmt.Println(query)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"employee_role_id": params.EmployeeRoleID,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeRolePermissionStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.EmployeeRolePermission{
			EmployeeRoleID: v.EmployeeRoleID,
			PermissionID:   v.PermissionID,
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

func (s EmployeeRolePermissionRepository) Find(ctx *gin.Context, id string) (*models.EmployeeRolePermission, *types.Error) {
	result := &models.EmployeeRolePermission{}
	bulks := []*models.EmployeeRolePermissionBulk{}

	var err error

	query := `
  SELECT
    employee_role_permissions.employee_role_id,
    employee_role_permissions.permission_id,
    permissions.package AS permission_package,
    permissions.module_name AS permission_module_name,
    permissions.action_name AS permission_action_name,
    permissions.display_module_name AS permission_display_module_name,
    permissions.display_action_name AS permission_display_action_name,
    permissions.http_method AS permission_http_method,
    permissions.route AS permission_route
  FROM employee_role_permissions
  JOIN permissions on permissions.id = employee_role_permissions.permission_id
  WHERE employee_role_permissions.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{})
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeRolePermissionStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = &models.EmployeeRolePermission{
			EmployeeRoleID: v.EmployeeRoleID,
			PermissionID:   v.PermissionID,
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
			Path:       ".EmployeeRoleStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return result, nil
}

func (s EmployeeRolePermissionRepository) Create(ctx *gin.Context, obj *models.CreateUpdateEmployeeRolePermission) (*models.EmployeeRolePermission, *types.Error) {
	data := models.EmployeeRolePermission{}
	_, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".EmployeeRolePermissionStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	data.EmployeeRoleID = obj.EmployeeRoleID
	data.PermissionID = obj.PermissionID
	return &data, nil
}

func (s EmployeeRolePermissionRepository) DeleteByEmployeeRoleID(ctx *gin.Context, id string) *types.Error {
	args := make(map[string]interface{})
	err := s.repository.ExecQuery(ctx, fmt.Sprintf(`DELETE FROM employee_role_permissions WHERE employee_role_id = "%s"`, id), args)
	if err != nil {
		return &types.Error{
			Path:       ".EmployeeRolePermissionStorage->DeleteByEmployeeRoleID()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return nil
}

func (s EmployeeRolePermissionRepository) CreateBunch(ctx *gin.Context, userID string, params models.FindAllEmployeeRolePermissionParams) *types.Error {
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
  INSERT INTO employee_role_permissions (employee_role_id, permission_id, created_at, updated_at)
  SELECT "%s", id, UTC_TIMESTAMP + INTERVAL 7 hour, UTC_TIMESTAMP + INTERVAL 7 HOUR
  FROM (
    SELECT id FROM permissions
    WHERE %s AND id %s IN (
      SELECT permission_id FROM employee_role_permissions
      WHERE employee_role_id = "%s"
    )
  ) permission`, userID, where, not, params.EmployeeRoleID)

	err := s.repository.ExecQuery(ctx, query, args)
	if err != nil {
		return &types.Error{
			Path:       ".EmployeeRolePermissionStorage->CreateBunch()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return nil
}
