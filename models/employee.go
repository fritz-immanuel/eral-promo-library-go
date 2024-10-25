package models

import "github.com/fritz-immanuel/eral-promo-library-go/library/types"

type EmployeeBulk struct {
	ID       string `json:"ID" db:"id"`
	Name     string `json:"Name" db:"name"`
	Email    string `json:"Email" db:"email"`
	Username string `json:"Username" db:"username"`
	Password string `json:"Password" db:"password"`
	StatusID string `json:"StatusID" db:"status_id"`

	StatusName string `json:"StatusName" db:"status_name"`
}

type Employee struct {
	ID       string `json:"ID" db:"id"`
	Name     string `json:"Name" db:"name" validate:"required"`
	Email    string `json:"Email" db:"email"`
	Username string `json:"Username" db:"username" validate:"required"`
	Password string `json:"Password" db:"password" validate:"required"`
	StatusID string `json:"StatusID" db:"status_id"`

	Status Status `json:"Status"`

	Permission []*Permission `json:"Permission"`
}

type FindAllEmployeeParams struct {
	FindAllParams types.FindAllParams
	Email         string
	Username      string
	Password      string
}

type EmployeeLogin struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	Password string `json:"Password" validate:"required"`
	Token    string `json:"Token"`
	Username string `json:"Username" validate:"required"`
}
