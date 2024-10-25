package models

import (
	"time"

	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
)

type Permission struct {
	ID                int    `json:"ID" db:"id"`
	Package           string `json:"Package" db:"package"`
	ModuleName        string `json:"ModuleName" db:"module_name"`
	ActionName        string `json:"ActionName" db:"action_name"`
	DisplayModuleName string `json:"DisplayModuleName" db:"display_module_name"`
	DisplayActionName string `json:"DisplayActionName" db:"display_action_name"`
	HTTPMethod        string `json:"HTTPMethod" db:"http_method"`
	Route             string `json:"Route" db:"route"`
	TableName         string `json:"TableName" db:"table_name"`

	CreatedAt time.Time `json:"CreatedAt" db:"created_at"`
	CreatedBy int       `json:"CreatedBy" db:"created_by"`
	UpdatedAt time.Time `json:"UpdatedAt" db:"updated_at"`
	UpdatedBy int       `json:"UpdatedBy" db:"updated_by"`
}

type FindAllPermissionParams struct {
	FindAllParams types.FindAllParams
	Package       string
	Name          string
	IsHidden      int
}
