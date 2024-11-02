package models

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
)

type EmployeeRolePermissionBulk struct {
	ID             string `json:"ID" db:"id"`
	EmployeeRoleID string `json:"EmployeeRoleID" db:"employee_role_id"`
	PermissionID   int    `json:"PermissionID" db:"permission_id"`

	PermissionPackage           string `json:"PermissionPackage" db:"permission_package"`
	PermissionModuleName        string `json:"PermissionModuleName" db:"permission_module_name"`
	PermissionActionName        string `json:"PermissionActionName" db:"permission_action_name"`
	PermissionDisplayModuleName string `json:"PermissionDisplayModuleName" db:"permission_display_module_name"`
	PermissionDisplayActionName string `json:"PermissionDisplayActionName" db:"permission_display_action_name"`
	PermissionHTTPMethod        string `json:"PermissionHTTPMethod" db:"permission_http_method"`
	PermissionRoute             string `json:"PermissionRoute" db:"permission_route"`

	Permission Permission `json:"Permission" db:"-"`
}

type EmployeeRolePermission struct {
	EmployeeRoleID string     `json:"EmployeeRoleID" db:"employee_role_id"`
	PermissionID   int        `json:"PermissionID" db:"permission_id"`
	Permission     Permission `json:"Permission" db:"-"`
}

type CreateUpdateEmployeeRolePermission struct {
	ID             string `json:"ID" db:"id"`
	EmployeeRoleID string `json:"EmployeeRoleID" db:"employee_role_id"`
	PermissionID   int    `json:"PermissionID" db:"permission_id"`
}

type FindAllEmployeeRolePermissionParams struct {
	FindAllParams      types.FindAllParams
	Package            string
	PermissionIDString string
	Not                int
	EmployeeRoleID     string
}
