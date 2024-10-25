package models

import "github.com/fritz-immanuel/eral-promo-library-go/library/types"

type UserBulk struct {
	ID         string `json:"ID" db:"id"`
	Name       string `json:"Name" db:"name"`
	Email      string `json:"Email" db:"email"`
	Username   string `json:"Username" db:"username"`
	Password   string `json:"Password" db:"password"`
	BusinessID string `json:"BusinessID" db:"business_id"`
	StatusID   string `json:"StatusID" db:"status_id"`

	StatusName string `json:"StatusName" db:"status_name"`
}

type User struct {
	ID         string `json:"ID" db:"id"`
	Name       string `json:"Name" db:"name" validate:"required"`
	Email      string `json:"Email" db:"email"`
	Username   string `json:"Username" db:"username" validate:"required"`
	Password   string `json:"Password" db:"password" validate:"required"`
	BusinessID string `json:"BusinessID" db:"business_id" validate:"required"`
	StatusID   string `json:"StatusID" db:"status_id"`

	Status Status `json:"Status"`

	Permission []*Permission `json:"Permission"`
}

type FindAllUserParams struct {
	FindAllParams types.FindAllParams
	Email         string
	Username      string
	Password      string
	BusinessID    string
}

type UserLogin struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	Password string `json:"Password" validate:"required"`
	Token    string `json:"Token"`
	Username string `json:"Username" validate:"required"`
}

type UserLoginAPI struct {
	Result UserLoginData `json:"result"`
}

type UserLoginData struct {
	Data *UserLogin `json:"Data"`
}
