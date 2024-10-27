package usecase

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/employeerole"
	"github.com/google/uuid"

	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type EmployeeRoleUsecase struct {
	employeeroleRepo           employeerole.Repository
	employeerolepermissionRepo employeerole.PermissionRepository
	contextTimeout             time.Duration
	db                         *sqlx.DB
}

func NewEmployeeRoleUsecase(db *sqlx.DB, employeeroleRepo employeerole.Repository, employeerolepermissionRepo employeerole.PermissionRepository) employeerole.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &EmployeeRoleUsecase{
		employeeroleRepo:           employeeroleRepo,
		employeerolepermissionRepo: employeerolepermissionRepo,
		contextTimeout:             timeoutContext,
		db:                         db,
	}
}

func (u *EmployeeRoleUsecase) FindAll(ctx *gin.Context, params models.FindAllEmployeeRoleParams) ([]*models.EmployeeRole, *types.Error) {
	result, err := u.employeeroleRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".EmployeeRoleUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *EmployeeRoleUsecase) Find(ctx *gin.Context, id string) (*models.EmployeeRole, *types.Error) {
	result, err := u.employeeroleRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".EmployeeRoleUsecase->Find()" + err.Path
		return nil, err
	}

	result.Permission, err = u.employeerolepermissionRepo.FindAll(ctx, models.FindAllEmployeeRolePermissionParams{EmployeeRoleID: id})
	if err != nil {
		err.Path = ".EmployeeRoleUsecase->Find()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *EmployeeRoleUsecase) Count(ctx *gin.Context, params models.FindAllEmployeeRoleParams) (int, *types.Error) {
	result, err := u.employeeroleRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".EmployeeRoleUsecase->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}

func (u *EmployeeRoleUsecase) Create(ctx *gin.Context, obj models.EmployeeRole) (*models.EmployeeRole, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".EmployeeRoleUsecase->Create()" + err.Path
		return nil, err
	}

	data := models.EmployeeRole{}
	data.ID = uuid.New().String()
	data.Name = obj.Name
	data.IsSupervisor = obj.IsSupervisor
	data.StatusID = models.DEFAULT_STATUS_CODE

	result, err := u.employeeroleRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".EmployeeRoleUsecase->Create()" + err.Path
		return nil, err
	}

	// create permission
	var permssions []string
	for _, v := range obj.Permission {
		permssions = append(permssions, fmt.Sprintf(`%d`, v.PermissionID))
	}

	var permissionParams models.FindAllEmployeeRolePermissionParams
	permissionParams.PermissionIDString = strings.Join(permssions, ",")
	err = u.employeerolepermissionRepo.CreateBunch(ctx, data.ID, permissionParams)
	if err != nil {
		err.Path = ".EmployeeRoleUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *EmployeeRoleUsecase) Update(ctx *gin.Context, id string, obj models.EmployeeRole) (*models.EmployeeRole, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".EmployeeRoleUsecase->Update()" + err.Path
		return nil, err
	}

	data, err := u.employeeroleRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".EmployeeRoleUsecase->Update()" + err.Path
		return nil, err
	}

	data.Name = obj.Name
	data.IsSupervisor = obj.IsSupervisor

	result, err := u.employeeroleRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".EmployeeRoleUsecase->Update()" + err.Path
		return nil, err
	}

	// update permission
	err = u.employeerolepermissionRepo.DeleteByEmployeeRoleID(ctx, id)
	if err != nil {
		err.Path = ".EmployeeRoleUsecase->Update()" + err.Path
		return nil, err
	}

	var permssions []string
	for _, v := range obj.Permission {
		permssions = append(permssions, fmt.Sprintf(`%d`, v.PermissionID))
	}

	var permissionParams models.FindAllEmployeeRolePermissionParams
	permissionParams.PermissionIDString = strings.Join(permssions, ",")
	err = u.employeerolepermissionRepo.CreateBunch(ctx, data.ID, permissionParams)
	if err != nil {
		err.Path = ".EmployeeRoleUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *EmployeeRoleUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.EmployeeRole, *types.Error) {
	if newStatusID != models.STATUS_ACTIVE && newStatusID != models.STATUS_INACTIVE {
		return nil, &types.Error{
			Path:       ".EmployeeRoleUsecase->UpdateStatus()",
			Message:    "StatusID is not valid",
			Error:      fmt.Errorf("StatusID is not valid"),
			StatusCode: http.StatusBadRequest,
		}
	}

	result, err := u.employeeroleRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".EmployeeRoleUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, err
}
