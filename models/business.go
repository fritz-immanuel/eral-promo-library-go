package models

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
)

type BusinessBulk struct {
	ID       string `json:"ID" db:"id"`
	Name     string `json:"Name" db:"name"`
	Code     string `json:"Code" db:"code"`
	LogoImg  string `json:"LogoImg" db:"logo_img"`
	StatusID string `json:"StatusID" db:"status_id"`

	StatusName string `json:"StatusName" db:"status_name"`
}

type Business struct {
	ID       string `json:"ID" db:"id"`
	Name     string `json:"Name" db:"name"`
	Code     string `json:"Code" db:"code"`
	LogoImg  string `json:"LogoImg" db:"logo_img"`
	StatusID string `json:"StatusID" db:"status_id"`

	Status Status `json:"Status"`
}

type FindAllBusinessParams struct {
	FindAllParams types.FindAllParams
}
