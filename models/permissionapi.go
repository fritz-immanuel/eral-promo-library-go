package models

type PermissionAPIResult struct {
	ID                uint   `json:"ID" db:"id"`
	Package           string `json:"Package" db:"package"`
	Name              string `json:"Name" db:"name"`
	Action            string `json:"Action" db:"action"`
	Type              string `json:"Type" db:"type"`
	Route             string `json:"Route" db:"route"`
	IsHidden          int    `json:"IsHidden" db:"is_hidden"`
	TableName         string `json:"TableName" db:"table_name"`
	DisplayActionName string `json:"DisplayActionName" db:"display_action_name"`
	DisplayModuleName string `json:"DisplayModuleName" db:"display_module_name"`
}
