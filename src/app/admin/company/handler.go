package company

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/eral-promo-library-go/library/firebase"
	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/middleware"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/company"
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/http/response"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"

	companyRepository "github.com/fritz-immanuel/eral-promo-library-go/src/services/company/repository"
	companyUsecase "github.com/fritz-immanuel/eral-promo-library-go/src/services/company/usecase"
)

type CompanyHandler struct {
	CompanyUsecase company.Usecase
	dataManager    *data.Manager
	Result         gin.H
	Status         int
}

func (h CompanyHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	companyRepo := companyRepository.NewCompanyRepository(
		data.NewMySQLStorage(db, "companies", models.Company{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	uCompany := companyUsecase.NewCompanyUsecase(db, companyRepo)

	base := &CompanyHandler{CompanyUsecase: uCompany, dataManager: dataManager}

	rs := v.Group("/company")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("/:id", middleware.Auth, base.Update)

		rs.PUT("/:id/status", middleware.Auth, base.UpdateStatus)
	}

	rss := v.Group("/statuses")
	{
		rss.GET("/company", base.FindStatus)
	}
}

func (h *CompanyHandler) FindAll(c *gin.Context) {
	var params models.FindAllCompanyParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	params.FindAllParams.SortBy = "companies.name ASC"
	datas, err := h.CompanyUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}
	}

	for _, data := range datas {
		if data.LogoImgURL != "" {
			data.LogoImgURL, _ = firebase.GenerateSignedURL(data.LogoImgURL)
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.CompanyUsecase.Count(c, params)
	if err != nil {
		err.Path = ".CompanyHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "Company Data fetched!", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *CompanyHandler) Find(c *gin.Context) {
	id := c.Param("id")

	result, err := h.CompanyUsecase.Find(c, id)
	if err != nil {
		err.Path = ".CompanyHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "Company not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	if result.LogoImgURL != "" {
		result.LogoImgURL, _ = firebase.GenerateSignedURL(result.LogoImgURL)
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Company Data fetched!", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *CompanyHandler) Create(c *gin.Context) {
	var err *types.Error
	var company models.Company
	var dataCompany *models.Company

	company.Name = c.PostForm("Name")
	company.Code = c.PostForm("Code")

	file, errFile := c.FormFile("LogoImgURL")
	if file != nil {
		if errFile != nil {
			err = &types.Error{
				Path:       ".CompanyHandler->Create()",
				Message:    errFile.Error(),
				Error:      errFile,
				StatusCode: http.StatusInternalServerError,
				Type:       "golang-error",
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		filename, err := firebase.UploadFile(c, file, "company")
		if err != nil {
			err.Path = ".CompanyHandler->Create()" + err.Path
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		company.LogoImgURL = filename
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		dataCompany, err = h.CompanyUsecase.Create(c, company)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".CompanyHandler->Create()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Company Data created!", Data: dataCompany}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *CompanyHandler) Update(c *gin.Context) {
	var err *types.Error
	var company models.Company
	var data *models.Company

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".CompanyHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	company.Name = c.PostForm("Name")
	company.Code = c.PostForm("Code")

	file, errFile := c.FormFile("LogoImgURL")
	if file != nil {
		if errFile != nil {
			err = &types.Error{
				Path:       ".CompanyHandler->Update()",
				Message:    errFile.Error(),
				Error:      errFile,
				StatusCode: http.StatusInternalServerError,
				Type:       "golang-error",
			}
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		filename, err := firebase.UploadFile(c, file, "company")
		if err != nil {
			err.Path = ".CompanyHandler->Update()" + err.Path
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}

		company.LogoImgURL = filename
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		{ // delete existing logo
			companyData, err := h.CompanyUsecase.Find(c, id)
			if err != nil {
				return err
			}

			if companyData.LogoImgURL != "" {
				err := firebase.DeleteFile(c, companyData.LogoImgURL)
				if err != nil {
					return err
				}
			}
		}

		data, err = h.CompanyUsecase.Update(c, id, company)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".CompanyHandler->Update()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Company Data updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *CompanyHandler) FindStatus(c *gin.Context) {
	var datas []*models.Status
	datas = append(datas, &models.Status{ID: models.STATUS_INACTIVE, Name: "Inactive"})
	datas = append(datas, &models.Status{ID: models.STATUS_ACTIVE, Name: "Active"})

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Company Status Data fetched!", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(http.StatusOK, h.Result)
}

func (h *CompanyHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.Company

	companyID, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".CompanyHandler->UpdateStatus()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	newStatusID := c.PostForm("StatusID")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.CompanyUsecase.UpdateStatus(c, companyID, newStatusID)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".CompanyHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Company Status has been updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
