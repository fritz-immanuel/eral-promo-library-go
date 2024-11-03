package models

import "github.com/fritz-immanuel/eral-promo-library-go/library/types"

type EmployeeBrandBulk struct {
	ID         string `json:"ID" db:"id"`
	EmployeeID string `json:"EmployeeID" db:"employee_id"`
	BrandID    string `json:"BrandID" db:"brand_id"`

	BrandName string `json:"BrandName" db:"brand_name"`
	BrandCode string `json:"BrandCode" db:"brand_code"`
}

type EmployeeBrand struct {
	ID         string `json:"ID" db:"id"`
	EmployeeID string `json:"EmployeeID" db:"employee_id" validate:"required"`
	BrandID    string `json:"BrandID" db:"brand_id" validate:"required"`

	Brand *StringIDNameCodeTemplate `json:"Brand"`
}

type FindAllEmployeeBrandParams struct {
	FindAllParams types.FindAllParams
	EmployeeID    string
	BrandID       string
}
