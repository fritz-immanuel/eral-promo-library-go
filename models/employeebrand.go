package models

import "github.com/fritz-immanuel/eral-promo-library-go/library/types"

type EmployeeBrandBulk struct {
	ID         string `json:"ID" db:"id"`
	EmployeeID string `json:"EmployeeID" db:"employee_id"`
	BrandID    string `json:"BrandID" db:"brand_id"`
}

type EmployeeBrand struct {
	ID         string `json:"ID" db:"id"`
	EmployeeID string `json:"EmployeeID" db:"name" validate:"required"`
	BrandID    string `json:"BrandID" db:"brand_id" validate:"required"`
}

type FindAllEmployeeBrandParams struct {
	FindAllParams types.FindAllParams
	EmployeeID    string
	BrandID       string
}
