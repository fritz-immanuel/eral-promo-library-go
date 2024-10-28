package models

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
)

type BrandBulk struct {
	ID         string `json:"ID" db:"id"`
	Name       string `json:"Name" db:"name"`
	Code       string `json:"Code" db:"code"`
	LogoImgURL string `json:"LogoImgURL" db:"logo_img_url"`
	BusinessID string `json:"BusinessID" db:"business_id"`
	StatusID   string `json:"StatusID" db:"status_id"`

	StatusName string `json:"StatusName" db:"status_name"`
}

type Brand struct {
	ID         string `json:"ID" db:"id"`
	Name       string `json:"Name" db:"name"`
	Code       string `json:"Code" db:"code"`
	LogoImgURL string `json:"LogoImgURL" db:"logo_img_url"`
	BusinessID string `json:"BusinessID" db:"business_id"`
	StatusID   string `json:"StatusID" db:"status_id"`

	Status Status `json:"Status"`
}

type FindAllBrandParams struct {
	FindAllParams types.FindAllParams
	BusinessID    string
}
