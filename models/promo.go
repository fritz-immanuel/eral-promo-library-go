package models

import (
	"time"

	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
)

type PromoBulk struct {
	ID               string     `json:"ID" db:"id"`
	Name             string     `json:"Name" db:"name"`
	Code             string     `json:"Code" db:"code"`
	StartDate        time.Time  `json:"StartDate" db:"start_date"`
	EndDate          time.Time  `json:"EndDate" db:"end_date"`
	ImgURL           string     `json:"ImgURL" db:"img_url"`
	CompanyID        string     `json:"CompanyID" db:"company_id"`
	BusinessID       string     `json:"BusinessID" db:"business_id"`
	BrandID          string     `json:"BrandID" db:"brand_id"`
	TotalPromoBudget float64    `json:"TotalPromoBudget" db:"total_promo_budget"`
	PrincipleSupport float64    `json:"PrincipleSupport" db:"principle_support"`
	InternalSupport  float64    `json:"InternalSupport" db:"internal_support"`
	Description      string     `json:"Description" db:"description"`
	ApprovedAt       *time.Time `json:"ApprovedAt" db:"approved_at"`
	ApprovedBy       string     `json:"ApprovedBy" db:"approved_by"`
	RejectedAt       *time.Time `json:"RejectedAt" db:"rejected_at"`
	RejectedBy       string     `json:"RejectedBy" db:"rejected_by"`
	RejectReason     string     `json:"RejectReason" db:"reject_reason"`
	StatusID         string     `json:"StatusID" db:"status_id"`

	StatusName string `json:"StatusName" db:"status_name"`
}

type Promo struct {
	ID               string     `json:"ID" db:"id"`
	Name             string     `json:"Name" db:"name" validate:"required"`
	Code             string     `json:"Code" db:"code"`
	StartDate        time.Time  `json:"StartDate" db:"start_date" validate:"required"`
	EndDate          time.Time  `json:"EndDate" db:"end_date" validate:"required"`
	ImgURL           string     `json:"ImgURL" db:"img_url"`
	CompanyID        string     `json:"CompanyID" db:"company_id"`
	BusinessID       string     `json:"BusinessID" db:"business_id"`
	BrandID          string     `json:"BrandID" db:"brand_id" validate:"required"`
	TotalPromoBudget float64    `json:"TotalPromoBudget" db:"total_promo_budget" validate:"required"`
	PrincipleSupport float64    `json:"PrincipleSupport" db:"principle_support" validate:"required"`
	InternalSupport  float64    `json:"InternalSupport" db:"internal_support" validate:"required"`
	Description      string     `json:"Description" db:"description"`
	ApprovedAt       *time.Time `json:"ApprovedAt" db:"approved_at"`
	ApprovedBy       string     `json:"ApprovedBy" db:"approved_by"`
	RejectedAt       *time.Time `json:"RejectedAt" db:"rejected_at"`
	RejectedBy       string     `json:"RejectedBy" db:"rejected_by"`
	RejectReason     string     `json:"RejectReason" db:"reject_reason"`
	StatusID         string     `json:"StatusID" db:"status_id"`

	Status Status `json:"Status"`

	PromoDocuments []*PromoDocument `json:"PromoDocuments"`
}

type FindAllPromoParams struct {
	FindAllParams  types.FindAllParams
	PromoTypeID    string
	CompanyID      string
	BusinessID     string
	BrandID        string
	Code           string
	StartDate      *time.Time
	EndDate        *time.Time
	ApprovalStatus int
}
