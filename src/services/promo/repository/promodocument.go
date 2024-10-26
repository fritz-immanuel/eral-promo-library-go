package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

type PromoDocumentRepository struct {
	repository data.GenericStorage
}

func NewPromoDocumentRepository(repository data.GenericStorage) PromoDocumentRepository {
	return PromoDocumentRepository{repository: repository}
}

// A function to get all Data that matches the filter provided
func (s PromoDocumentRepository) FindAll(ctx *gin.Context, params models.FindAllPromoDocumentParams) ([]*models.PromoDocument, *types.Error) {
	result := []*models.PromoDocument{}
	bulks := []*models.PromoDocumentBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where += fmt.Sprintf(` AND promo_documents.%s`, params.FindAllParams.StatusID)
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT
    promo_documents.id, promo_documents.promo_id, promo_documents.document_url,
    promo_documents.status_id
  FROM promo_documents
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":  params.FindAllParams.Size,
		"offset": ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoDocumentStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		for _, v := range bulks {
			result = append(result, &models.PromoDocument{
				ID:          v.ID,
				PromoID:     v.PromoID,
				DocumentURL: v.DocumentURL,
				StatusID:    v.StatusID,
			})
		}
	}

	return result, nil
}

// A function to get a row of data specified by the given ID
func (s PromoDocumentRepository) Find(ctx *gin.Context, id string) (*models.PromoDocument, *types.Error) {
	result := models.PromoDocument{}
	bulks := []*models.PromoDocumentBulk{}
	var err error

	query := `
  SELECT
    promo_documents.id, promo_documents.promo_id, promo_documents.document_url,
    promo_documents.status_id
  FROM promo_documents
  WHERE promo_documents.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoDocumentStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.PromoDocument{
			ID:          v.ID,
			PromoID:     v.PromoID,
			DocumentURL: v.DocumentURL,
			StatusID:    v.StatusID,
		}
	} else {
		return nil, &types.Error{
			Path:       ".PromoDocumentStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

// Inserts a new row of data
func (s PromoDocumentRepository) Create(ctx *gin.Context, obj *models.PromoDocument) (*models.PromoDocument, *types.Error) {
	result := models.PromoDocument{}
	_, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoDocumentStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &result, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoDocumentStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

// Updates a row of data specified by the given ID inside the obj struct
func (s PromoDocumentRepository) Update(ctx *gin.Context, obj *models.PromoDocument) (*models.PromoDocument, *types.Error) {
	result := models.PromoDocument{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoDocumentStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &result, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoDocumentStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s PromoDocumentRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.PromoDocument, *types.Error) {
	data := models.PromoDocument{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoDocumentStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".PromoDocumentStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}
