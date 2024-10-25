package usecase

import (
	"time"

	"github.com/fritz-immanuel/eral-promo-library-go/library/helpers"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/userpermission"
	"github.com/google/uuid"

	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type UserPermissionUsecase struct {
	userpermissionRepo userpermission.Repository
	contextTimeout     time.Duration
	db                 *sqlx.DB
}

func NewUserPermissionUsecase(db *sqlx.DB, userpermissionRepo userpermission.Repository) userpermission.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	return &UserPermissionUsecase{
		userpermissionRepo: userpermissionRepo,
		contextTimeout:     timeoutContext,
		db:                 db,
	}
}

func (u *UserPermissionUsecase) FindAll(ctx *gin.Context, params models.FindAllUserPermissionParams) ([]*models.UserPermission, *types.Error) {
	result, err := u.userpermissionRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".UserPermissionUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *UserPermissionUsecase) Find(ctx *gin.Context, id string) (*models.UserPermission, *types.Error) {
	result, err := u.userpermissionRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".UserPermissionUsecase->Find()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *UserPermissionUsecase) Create(ctx *gin.Context, obj models.CreateUpdateUserPermission) (*models.UserPermission, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".UserPermissionUsecase->Create()" + err.Path
		return nil, err
	}

	var dataCreate models.CreateUpdateUserPermission
	dataCreate.ID = uuid.New().String()
	dataCreate.UserID = obj.UserID
	dataCreate.PermissionID = obj.PermissionID

	result, err := u.userpermissionRepo.Create(ctx, &dataCreate)
	if err != nil {
		err.Path = ".UserPermissionUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *UserPermissionUsecase) DeleteByUserID(ctx *gin.Context, id string) *types.Error {
	err := u.userpermissionRepo.DeleteByUserID(ctx, id)
	if err != nil {
		err.Path = ".UserPermissionUsecase->DeleteByUserID()" + err.Path
		return err
	}

	return nil
}

func (u *UserPermissionUsecase) CreateBunch(ctx *gin.Context, userID string, params models.FindAllUserPermissionParams) *types.Error {
	err := u.userpermissionRepo.CreateBunch(ctx, userID, params)
	if err != nil {
		err.Path = ".UserPermissionUsecase->CreateBunch()" + err.Path
		return err
	}

	return nil
}
