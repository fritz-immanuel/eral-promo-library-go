package models

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
)

type EmployeeRoleBulk struct {
	ID           string `json:"ID" db:"id"`
	Name         string `json:"Name" db:"name"`
	IsSupervisor int    `json:"IsSupervisor" db:"is_supervisor"`
	StatusID     string `json:"StatusID" db:"status_id"`

	StatusName string `json:"StatusName" db:"status_name"`
}

type EmployeeRole struct {
	ID           string `json:"ID" db:"id"`
	Name         string `json:"Name" db:"name"`
	IsSupervisor int    `json:"IsSupervisor" db:"is_supervisor"`
	StatusID     string `json:"StatusID" db:"status_id"`

	Status Status `json:"Status"`

	Permission []*EmployeeRolePermission `json:"Permission"`
}

type FindAllEmployeeRoleParams struct {
	FindAllParams   types.FindAllParams
	IsSupervisor    int
	IsNotSupervisor int
}
