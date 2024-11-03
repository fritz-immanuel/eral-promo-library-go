package promo

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/eral-promo-library-go/library"
	"github.com/fritz-immanuel/eral-promo-library-go/library/appcontext"
	"github.com/fritz-immanuel/eral-promo-library-go/library/firebase"
	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/middleware"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/promo"
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/http/response"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"

	promoRepository "github.com/fritz-immanuel/eral-promo-library-go/src/services/promo/repository"
	promoUsecase "github.com/fritz-immanuel/eral-promo-library-go/src/services/promo/usecase"
)

type PromoHandler struct {
	PromoUsecase promo.Usecase
	dataManager  *data.Manager
	Result       gin.H
	Status       int
}

func (h PromoHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	promoRepo := promoRepository.NewPromoRepository(
		data.NewMySQLStorage(db, "promos", models.Promo{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "promo_status", models.Status{}, data.MysqlConfig{}),
	)

	promodocumentRepo := promoRepository.NewPromoDocumentRepository(
		data.NewMySQLStorage(db, "promo_documents", models.PromoDocument{}, data.MysqlConfig{}),
	)

	uPromo := promoUsecase.NewPromoUsecase(db, promoRepo, promodocumentRepo)

	base := &PromoHandler{PromoUsecase: uPromo, dataManager: dataManager}

	rs := v.Group("/promos")
	{
		rs.GET("", middleware.AuthWebApp, base.FindAll)
		rs.GET("/:id", middleware.AuthWebApp, base.Find)
		rs.POST("", middleware.AuthWebApp, base.Create)
		rs.PUT("/:id", middleware.AuthWebApp, base.Update)

		rs.PUT("/:id/status", middleware.AuthWebApp, base.UpdateStatus)

		rs.POST("/:id/document", middleware.AuthWebApp, base.CreateDocument)
		rs.PUT("/:id/document/:documentID", middleware.AuthWebApp, base.UpdateDocument)
		rs.DELETE("/:id/document/:documentID", middleware.AuthWebApp, base.DeleteDocument)

		// Approval
		rs.PUT("/:id/approve", middleware.AuthWebApp, base.ApprovePromo)
		rs.PUT("/:id/reject", middleware.AuthWebApp, base.RejectPromo)
	}

	rss := v.Group("/statuses")
	{
		rss.GET("/promos", base.FindStatus)
	}
}

func (h *PromoHandler) FindAll(c *gin.Context) {
	var params models.FindAllPromoParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	params.CompanyID = *appcontext.CompanyID(c)
	params.BusinessID = *appcontext.BusinessID(c)
	params.BrandID, _ = helpers.MultiValueUUIDCheck(c.Query("BrandID"))
	params.ApprovalStatus, _ = strconv.Atoi(c.Query("ApprovalStatus"))

	if c.Query("StartDate") != "" {
		startDateTime, errConversion := time.Parse(library.DateStampFormat(), c.Query("StartDate"))
		if errConversion != nil {
			err := &types.Error{
				Path:       ".PromoHandler->FindAll()",
				Message:    "Incorrect Start Date Format",
				Error:      errConversion,
				Type:       "conversion-error",
				StatusCode: http.StatusBadRequest,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}
		params.StartDate = &startDateTime
	}

	if c.Query("EndDate") != "" {
		endDateTime, errConversion := time.Parse(library.DateStampFormat(), c.Query("EndDate"))
		if errConversion != nil {
			err := &types.Error{
				Path:       ".PromoHandler->FindAll()",
				Message:    "Incorrect End Date Format",
				Error:      errConversion,
				Type:       "conversion-error",
				StatusCode: http.StatusBadRequest,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}
		params.EndDate = &endDateTime
	}

	if c.Query("SortBy") != "" || c.Query("SortName") != "" {
		params.FindAllParams.SortBy = "promos.name ASC"
	}
	datas, err := h.PromoUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	for _, data := range datas {
		if data.ImgURL != "" {
			data.ImgURL, _ = firebase.GenerateSignedURL(data.ImgURL)
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.PromoUsecase.Count(c, params)
	if err != nil {
		err.Path = ".PromoHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "Promo Data fetched!", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *PromoHandler) Find(c *gin.Context) {
	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".PromoHandler->Find()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	result, err := h.PromoUsecase.Find(c, id)
	if err != nil {
		err.Path = ".PromoHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "Promo not found", http.StatusNotFound, *err)
			return
		}
		response.Error(c, err.Message, http.StatusInternalServerError, *err)
		return
	}

	if result.CompanyID != *appcontext.CompanyID(c) || result.BusinessID != *appcontext.BusinessID(c) {
		response.Error(c, "Promo not found", http.StatusNotFound, *err)
		return
	}

	if result.ImgURL != "" {
		result.ImgURL, _ = firebase.GenerateSignedURL(result.ImgURL)
	}

	if len(result.PromoDocuments) > 0 {
		for _, docs := range result.PromoDocuments {
			docs.DocumentURL, _ = firebase.GenerateSignedURL(docs.DocumentURL)
		}
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Promo Data fetched!", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *PromoHandler) Create(c *gin.Context) {
	var err *types.Error
	var promo models.Promo
	var dataPromo *models.Promo

	if c.PostForm("StartDate") != "" {
		startDateTime, errConversion := time.Parse(library.DateStampFormat(), c.PostForm("StartDate"))
		if errConversion != nil {
			err := &types.Error{
				Path:       ".PromoHandler->Create()",
				Message:    "Incorrect Start Date Format",
				Error:      errConversion,
				Type:       "conversion-error",
				StatusCode: http.StatusBadRequest,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}
		promo.StartDate = startDateTime
	}

	if c.PostForm("EndDate") != "" {
		endDateTime, errConversion := time.Parse(library.DateStampFormat(), c.PostForm("EndDate"))
		if errConversion != nil {
			err := &types.Error{
				Path:       ".PromoHandler->Create()",
				Message:    "Incorrect End Date Format",
				Error:      errConversion,
				Type:       "conversion-error",
				StatusCode: http.StatusBadRequest,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}
		promo.EndDate = endDateTime
	}

	promo.Name = c.PostForm("Name")
	promo.Code = c.PostForm("Code")
	promo.CompanyID = *appcontext.CompanyID(c)
	promo.BusinessID = *appcontext.BusinessID(c)
	promo.BrandID, err = helpers.ValidateUUID(c.PostForm("BrandID"))
	if err != nil {
		err.Path = ".PromoHandler->Create()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}
	promo.TotalPromoBudget, _ = strconv.ParseFloat(c.PostForm("TotalPromoBudget"), 64)
	promo.PrincipleSupport, _ = strconv.ParseFloat(c.PostForm("PrincipleSupport"), 64)
	promo.InternalSupport, _ = strconv.ParseFloat(c.PostForm("InternalSupport"), 64)
	promo.Description = c.PostForm("Description")

	{ // upload img
		file, errFile := c.FormFile("ImgURL")
		if file != nil {
			if errFile != nil {
				err = &types.Error{
					Path:       ".PromoHandler->Create()",
					Message:    errFile.Error(),
					Error:      errFile,
					StatusCode: http.StatusInternalServerError,
					Type:       "golang-error",
				}
				response.Error(c, err.Message, err.StatusCode, *err)
				return
			}

			filename, err := firebase.UploadFile(c, file, "promo")
			if err != nil {
				err.Path = ".PromoHandler->Create()" + err.Path
				response.Error(c, err.Message, err.StatusCode, *err)
				return
			}

			promo.ImgURL = filename
		}
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		dataPromo, err = h.PromoUsecase.Create(c, promo)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".PromoHandler->Create()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Promo Data created!", Data: dataPromo}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *PromoHandler) Update(c *gin.Context) {
	var err *types.Error
	var promo models.Promo
	var data *models.Promo

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".PromoHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	if c.PostForm("StartDate") != "" {
		startDateTime, errConversion := time.Parse(library.DateStampFormat(), c.PostForm("StartDate"))
		if errConversion != nil {
			err := &types.Error{
				Path:       ".PromoHandler->Update()",
				Message:    "Incorrect Start Date Format",
				Error:      errConversion,
				Type:       "conversion-error",
				StatusCode: http.StatusBadRequest,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}
		promo.StartDate = startDateTime
	}

	if c.PostForm("EndDate") != "" {
		endDateTime, errConversion := time.Parse(library.DateStampFormat(), c.PostForm("EndDate"))
		if errConversion != nil {
			err := &types.Error{
				Path:       ".PromoHandler->Update()",
				Message:    "Incorrect End Date Format",
				Error:      errConversion,
				Type:       "conversion-error",
				StatusCode: http.StatusBadRequest,
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}
		promo.EndDate = endDateTime
	}

	promo.Name = c.PostForm("Name")
	promo.Code = c.PostForm("Code")
	promo.BrandID, err = helpers.ValidateUUID(c.PostForm("BrandID"))
	if err != nil {
		err.Path = ".PromoHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}
	promo.TotalPromoBudget, _ = strconv.ParseFloat(c.PostForm("TotalPromoBudget"), 64)
	promo.PrincipleSupport, _ = strconv.ParseFloat(c.PostForm("PrincipleSupport"), 64)
	promo.InternalSupport, _ = strconv.ParseFloat(c.PostForm("InternalSupport"), 64)
	promo.Description = c.PostForm("Description")

	{ // upload img
		file, errFile := c.FormFile("ImgURL")
		if file != nil {
			if errFile != nil {
				err = &types.Error{
					Path:       ".PromoHandler->Update()",
					Message:    errFile.Error(),
					Error:      errFile,
					StatusCode: http.StatusInternalServerError,
					Type:       "golang-error",
				}
				response.Error(c, err.Message, err.StatusCode, *err)
				return
			}

			filename, err := firebase.UploadFile(c, file, "promo")
			if err != nil {
				err.Path = ".PromoHandler->Update()" + err.Path
				response.Error(c, err.Message, err.StatusCode, *err)
				return
			}

			promo.ImgURL = filename
		}
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		{ // delete existing img
			promoData, err := h.PromoUsecase.Find(c, id)
			if err != nil {
				return err
			}

			if promoData.ImgURL != "" {
				err := firebase.DeleteFile(c, promoData.ImgURL)
				if err != nil {
					return err
				}
			}
		}

		data, err = h.PromoUsecase.Update(c, id, promo)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".PromoHandler->Update()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Promo Data updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *PromoHandler) FindStatus(c *gin.Context) {
	datas, err := h.PromoUsecase.FindStatus(c)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Promo Status fetched!", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(http.StatusOK, h.Result)
}

func (h *PromoHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.Promo

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".PromoHandler->UpdateStatus()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	newStatusID := c.PostForm("StatusID")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.PromoUsecase.UpdateStatus(c, id, newStatusID)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".PromoHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Promo Status has been updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

// APPROVAL

func (h *PromoHandler) ApprovePromo(c *gin.Context) {
	var err *types.Error
	var data *models.Promo

	if appcontext.IsSupervisor(c) == 0 {
		err = &types.Error{
			Path:       ".PromoHandler->ApprovePromo()",
			Message:    "You are not allowed to perform this action",
			Error:      fmt.Errorf("You are not allowed to perform this action"),
			StatusCode: http.StatusForbidden,
			Type:       "validation-error",
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".PromoHandler->ApprovePromo()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.PromoUsecase.ApprovePromo(c, id)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".PromoHandler->ApprovePromo()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Promo approved!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *PromoHandler) RejectPromo(c *gin.Context) {
	var err *types.Error
	var data *models.Promo

	if appcontext.IsSupervisor(c) == 0 {
		err = &types.Error{
			Path:       ".PromoHandler->RejectPromo()",
			Message:    "You are not allowed to perform this action",
			Error:      fmt.Errorf("You are not allowed to perform this action"),
			StatusCode: http.StatusForbidden,
			Type:       "validation-error",
		}
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".PromoHandler->RejectPromo()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	rejectReason := c.PostForm("RejectReason")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.PromoUsecase.RejectPromo(c, id, rejectReason)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".PromoHandler->RejectPromo()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Promo rejected!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

// DOCUMENT

func (h *PromoHandler) CreateDocument(c *gin.Context) {
	var err *types.Error
	var promoDoc models.PromoDocument
	var dataPromoDocument *models.PromoDocument

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".PromoHandler->CreateDocument()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	promoDoc.PromoID = id

	{ // upload document
		file, errFile := c.FormFile("Document")
		if file != nil {
			if errFile != nil {
				err = &types.Error{
					Path:       ".PromoHandler->CreateDocument()",
					Message:    errFile.Error(),
					Error:      errFile,
					StatusCode: http.StatusInternalServerError,
					Type:       "golang-error",
				}
				response.Error(c, err.Message, err.StatusCode, *err)
				return
			}

			filename, err := firebase.UploadFile(c, file, "promo")
			if err != nil {
				err.Path = ".PromoHandler->CreateDocument()" + err.Path
				response.Error(c, err.Message, err.StatusCode, *err)
				return
			}

			promoDoc.DocumentURL = filename
		}
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		dataPromoDocument, err = h.PromoUsecase.CreateDocument(c, promoDoc)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".PromoHandler->CreateDocument()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Promo Document Data created!", Data: dataPromoDocument}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *PromoHandler) UpdateDocument(c *gin.Context) {
	var err *types.Error
	var promoDoc models.PromoDocument
	var data *models.PromoDocument

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".PromoHandler->UpdateDocument()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	documentID, err := helpers.ValidateUUID(c.Param("documentID"))
	if err != nil {
		err.Path = ".PromoHandler->UpdateDocument()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	promoDoc.PromoID = id

	{ // upload document
		file, errFile := c.FormFile("Document")
		if file != nil {
			if errFile != nil {
				err = &types.Error{
					Path:       ".PromoHandler->UpdateDocument()",
					Message:    errFile.Error(),
					Error:      errFile,
					StatusCode: http.StatusInternalServerError,
					Type:       "golang-error",
				}
				response.Error(c, err.Message, err.StatusCode, *err)
				return
			}

			filename, err := firebase.UploadFile(c, file, "promo")
			if err != nil {
				err.Path = ".PromoHandler->UpdateDocument()" + err.Path
				response.Error(c, err.Message, err.StatusCode, *err)
				return
			}

			promoDoc.DocumentURL = filename
		}
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.PromoUsecase.UpdateDocument(c, documentID, promoDoc)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".PromoHandler->UpdateDocument()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Promo Document Data updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *PromoHandler) DeleteDocument(c *gin.Context) {
	var err *types.Error

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".PromoHandler->DeleteDocument()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	documentID, err := helpers.ValidateUUID(c.Param("documentID"))
	if err != nil {
		err.Path = ".PromoHandler->DeleteDocument()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		{ // check validity
			checkData, err := h.PromoUsecase.FindDocument(c, id)
			if err != nil {
				return err
			}

			if checkData.PromoID != id {
				err = &types.Error{
					Path:       ".PromoHandler->DeleteDocument()",
					Message:    "Data not found",
					Error:      fmt.Errorf("Promo ID does not match Document Promo ID"),
					StatusCode: http.StatusNotFound,
					Type:       "validation-error",
				}
				return err
			}
		}

		err = h.PromoUsecase.DeleteDocument(c, documentID)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".PromoHandler->DeleteDocument()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Promo Document Data deleted!", Data: nil}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

// HISTORY

func (h *PromoHandler) FindUserActionHistory(c *gin.Context) {
	var params models.FindAllActionHistory
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".PromoHandler->FindUserActionHistory()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	result, err := h.PromoUsecase.FindUserActionHistory(c, id, params)
	if err != nil {
		err.Path = ".PromoHandler->FindUserActionHistory()" + err.Path
		response.Error(c, err.Message, http.StatusInternalServerError, *err)
		return
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1

	leng, err := h.PromoUsecase.FindUserActionHistory(c, id, params)
	if err != nil {
		err.Path = ".PromoHandler->FindUserActionHistory()" + err.Path
		response.Error(c, err.Message, http.StatusInternalServerError, *err)
		return
	}

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "User Action Data fetched!", Data: result, Page: page, Size: size, TotalData: len(leng)}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
