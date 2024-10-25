package models

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
)

type BusinessConfigBulk struct {
	ID         uint   `json:"ID" db:"id"`
	BusinessID int    `json:"BusinessID" db:"business_id"`
	SubURLName string `json:"SubURLName" db:"sub_url_name"`
	Config     string `json:"Config" db:"config"`
}

type BusinessConfig struct {
	ID         uint   `json:"ID" db:"id"`
	BusinessID int    `json:"BusinessID" db:"business_id"`
	SubURLName string `json:"SubURLName" db:"sub_url_name"`
	Config     string `json:"Config" db:"config"`
}

type FindAllBusinessConfigParams struct {
	FindAllParams types.FindAllParams
	BusinessID    int
	SubURLName    string
}
