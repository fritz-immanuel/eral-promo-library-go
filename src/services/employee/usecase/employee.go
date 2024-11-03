package usecase

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fritz-immanuel/eral-promo-library-go/library"
	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/employee"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/employeerole"
	"github.com/google/uuid"

	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type EmployeeUsecase struct {
	employeeRepo        employee.Repository
	employeebrandRepo   employee.BrandRepository
	employeeroleUsecase employeerole.Usecase
	contextTimeout      time.Duration
	db                  *sqlx.DB
}

func NewEmployeeUsecase(db *sqlx.DB, employeeRepo employee.Repository, employeebrandRepo employee.BrandRepository, employeeroleUsecase employeerole.Usecase) employee.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &EmployeeUsecase{
		employeeRepo:        employeeRepo,
		employeebrandRepo:   employeebrandRepo,
		employeeroleUsecase: employeeroleUsecase,
		contextTimeout:      timeoutContext,
		db:                  db,
	}
}

func (u *EmployeeUsecase) FindAll(ctx *gin.Context, params models.FindAllEmployeeParams) ([]*models.Employee, *types.Error) {
	result, err := u.employeeRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".EmployeeUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *EmployeeUsecase) Find(ctx *gin.Context, id string) (*models.Employee, *types.Error) {
	result, err := u.employeeRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".EmployeeUsecase->Find()" + err.Path
		return nil, err
	}

	result.Brands, err = u.employeebrandRepo.FindAll(ctx, models.FindAllEmployeeBrandParams{EmployeeID: id})
	if err != nil {
		err.Path = ".EmployeeUsecase->Find()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *EmployeeUsecase) Count(ctx *gin.Context, params models.FindAllEmployeeParams) (int, *types.Error) {
	result, err := u.employeeRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".EmployeeUsecase->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}

func (u *EmployeeUsecase) Create(ctx *gin.Context, obj models.Employee) (*models.Employee, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".EmployeeUsecase->Create()" + err.Path
		return nil, err
	}

	// check for duplicate username
	employees, err := u.employeeRepo.FindAll(ctx, models.FindAllEmployeeParams{Username: obj.Username})
	if err != nil {
		err.Path = ".EmployeeUsecase->Create()" + err.Path
		return nil, err
	}

	if len(employees) > 0 {
		return nil, &types.Error{
			Path:       ".EmployeeUsecase->Create()",
			Message:    "Username already exists",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "mysql-error",
		}
	}

	data := models.Employee{}
	data.ID = uuid.New().String()
	data.Name = obj.Name
	data.Email = obj.Email
	data.Username = obj.Username
	data.Password = obj.Password
	data.BusinessID = obj.BusinessID
	data.EmployeeRoleID = obj.EmployeeRoleID
	data.StatusID = models.DEFAULT_STATUS_CODE

	result, err := u.employeeRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".EmployeeUsecase->Create()" + err.Path
		return nil, err
	}

	// Brands
	for _, v := range obj.Brands {
		v.EmployeeID = data.ID
		_, err := u.employeebrandRepo.Create(ctx, v)
		if err != nil {
			err.Path = ".EmployeeUsecase->Create()" + err.Path
			return nil, err
		}
	}

	return result, nil
}

func (u *EmployeeUsecase) Update(ctx *gin.Context, id string, obj models.Employee) (*models.Employee, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".EmployeeUsecase->Update()" + err.Path
		return nil, err
	}

	data, err := u.employeeRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".EmployeeUsecase->Update()" + err.Path
		return nil, err
	}

	// check for duplicate username
	var dupeParams models.FindAllEmployeeParams
	dupeParams.Username = obj.Username
	dupeParams.FindAllParams.DataFinder = fmt.Sprintf(`employees.id != '%s'`, id)
	employees, err := u.employeeRepo.FindAll(ctx, dupeParams)
	if err != nil {
		err.Path = ".EmployeeUsecase->Update()" + err.Path
		return nil, err
	}

	if len(employees) > 0 {
		return nil, &types.Error{
			Path:       ".EmployeeUsecase->Update()",
			Message:    "Username already exists",
			Error:      fmt.Errorf("Username already exists"),
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "mysql-error",
		}
	}

	data.Name = obj.Name
	data.Email = obj.Email
	data.Username = obj.Username
	data.BusinessID = obj.BusinessID
	data.EmployeeRoleID = obj.EmployeeRoleID

	result, err := u.employeeRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".EmployeeUsecase->Update()" + err.Path
		return nil, err
	}

	// Brands
	err = u.employeebrandRepo.DeleteByEmployeeID(ctx, data.ID)
	if err != nil {
		err.Path = ".EmployeeUsecase->Update()" + err.Path
		return nil, err
	}

	for _, v := range obj.Brands {
		v.EmployeeID = data.ID
		_, err := u.employeebrandRepo.Create(ctx, v)
		if err != nil {
			err.Path = ".EmployeeUsecase->Update()" + err.Path
			return nil, err
		}
	}

	return result, err
}

func (u *EmployeeUsecase) UpdatePassword(ctx *gin.Context, id string, newPassword string) (*models.Employee, *types.Error) {
	data, err := u.employeeRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".EmployeeUsecase->UpdatePassword()" + err.Path
		return nil, err
	}

	hash := md5.New()
	io.WriteString(hash, newPassword)
	data.Password = fmt.Sprintf("%x", hash.Sum(nil))

	result, err := u.employeeRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".EmployeeUsecase->UpdatePassword()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *EmployeeUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.Employee, *types.Error) {
	if newStatusID != models.STATUS_ACTIVE && newStatusID != models.STATUS_INACTIVE {
		return nil, &types.Error{
			Path:       ".EmployeeUsecase->UpdateStatus()",
			Message:    "StatusID is not valid",
			Error:      fmt.Errorf("StatusID is not valid"),
			StatusCode: http.StatusBadRequest,
		}
	}

	result, err := u.employeeRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".EmployeeUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *EmployeeUsecase) Login(ctx *gin.Context, creds models.EmployeeLogin) (*models.EmployeeLogin, *types.Error) {
	err := helpers.ValidateStruct(creds)
	if err != nil {
		err.Path = ".EmployeeUsecase->Login()" + err.Path
		return nil, err
	}

	var employeeParams models.FindAllEmployeeParams
	employeeParams.Username = creds.Username
	employeeParams.Password = creds.Password
	employeeParams.FindAllParams.StatusID = `employees.status_id = 1`
	employees, err := u.employeeRepo.FindAllForLogin(ctx, employeeParams)
	if err != nil {
		err.Path = ".EmployeeUsecase->Login()" + err.Path
		return nil, err
	}

	if len(employees) == 0 {
		return nil, &types.Error{
			StatusCode: http.StatusUnauthorized,
			Message:    "Username / Password is incorrect",
			Error:      data.ErrNotFound,
			Path:       ".EmployeeUsecase->Login()",
		}
	}

	employee := employees[0]

	credentials := library.CredentialWebApp{
		ID:             employee.ID,
		Name:           employee.Name,
		Email:          employee.Email,
		Username:       employee.Username,
		CompanyID:      employee.CompanyID,
		BusinessID:     employee.BusinessID,
		EmployeeRoleID: employee.EmployeeRoleID,
		IsSupervisor:   employee.IsSupervisor,
		Type:           "WebApp",
	}

	token, errorJwtSign := library.JwtSignWebAppString(credentials)
	if errorJwtSign != nil {
		return nil, &types.Error{
			Error:      errorJwtSign,
			Message:    "Error JWT Sign String",
			Path:       ".EmployeeUsecase->Login()",
			StatusCode: http.StatusInternalServerError,
		}
	}

	employeeRole, err := u.employeeroleUsecase.Find(ctx, employee.EmployeeRoleID)
	if err != nil {
		err.Path = ".EmployeeUsecase->Login()" + err.Path
		return nil, err
	}

	creds.Permissions = employeeRole.Permission
	creds.Name = employee.Name
	creds.Token = token
	creds.Password = ""

	return &creds, nil
}
