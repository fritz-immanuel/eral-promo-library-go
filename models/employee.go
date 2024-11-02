package models

import "github.com/fritz-immanuel/eral-promo-library-go/library/types"

type EmployeeBulk struct {
	ID             string `json:"ID" db:"id"`
	Name           string `json:"Name" db:"name"`
	Email          string `json:"Email" db:"email"`
	Username       string `json:"Username" db:"username"`
	Password       string `json:"Password" db:"password"`
	BusinessID     string `json:"BusinessID" db:"business_id"`
	EmployeeRoleID string `json:"EmployeeRoleID" db:"employee_role_id"`
	StatusID       string `json:"StatusID" db:"status_id"`

	StatusName string `json:"StatusName" db:"status_name"`
}

type Employee struct {
	ID             string `json:"ID" db:"id"`
	Name           string `json:"Name" db:"name" validate:"required"`
	Email          string `json:"Email" db:"email"`
	Username       string `json:"Username" db:"username" validate:"required"`
	Password       string `json:"Password" db:"password"`
	BusinessID     string `json:"BusinessID" db:"business_id" validate:"required"`
	EmployeeRoleID string `json:"EmployeeRoleID" db:"employee_role_id" validate:"required"`
	StatusID       string `json:"StatusID" db:"status_id"`

	Status Status `json:"Status"`

	Brands []*EmployeeBrand `json:"Brands"`
}

type FindAllEmployeeParams struct {
	FindAllParams  types.FindAllParams
	Email          string
	Username       string
	Password       string
	BusinessID     string
	EmployeeRoleID string
}

type EmployeeListForLogin struct {
	ID             string `json:"ID" db:"id"`
	Name           string `json:"Name" db:"name"`
	Email          string `json:"Email" db:"email"`
	Username       string `json:"Username" db:"username"`
	BusinessID     string `json:"BusinessID" db:"business_id"`
	EmployeeRoleID string `json:"EmployeeRoleID" db:"employee_role_id"`
	StatusID       string `json:"StatusID" db:"status_id"`

	CompanyID    string `json:"CompanyID" db:"company_id"`
	IsSupervisor int    `json:"IsSupervisor" db:"is_supervisor"`
}

type EmployeeLogin struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	Username string `json:"Username" validate:"required"`
	Password string `json:"Password" validate:"required"`
	Token    string `json:"Token"`

	Permissions []*EmployeeRolePermission `json:"Permissions"`
}
