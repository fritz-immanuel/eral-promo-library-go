package models

import (
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
)

type BusinessShift struct {
	ID         int    `json:"ID" db:"id"`
	BusinessID int    `json:"BusinessID" db:"business_id"`
	Name       string `json:"Name" db:"name"`
}

type FindAllBusinessShiftParams struct {
	FindAllParams types.FindAllParams
	BusinessID    int
}

// API
type BusinessShiftAllAPI struct {
	Result BusinessShiftData `json:"result" db:"-"`
}

type BusinessShiftData struct {
	TotalData int                 `json:"TotalData"`
	Data      []*BusinessShiftAPI `json:"Data" db:"-"`
}

type BusinessShiftAPI struct {
	ID         int    `json:"ID" db:"id"`
	BusinessID int    `json:"BusinessID" db:"business_id"`
	Name       string `json:"Name" db:"name"`
	DeletedBy  int    `json:"DeletedBy" db:"deleted_by"`
}
