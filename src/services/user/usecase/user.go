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
	"github.com/fritz-immanuel/eral-promo-library-go/src/services/user"
	"github.com/google/uuid"

	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type UserUsecase struct {
	userRepo       user.Repository
	contextTimeout time.Duration
	db             *sqlx.DB
}

func NewUserUsecase(db *sqlx.DB, userRepo user.Repository) user.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &UserUsecase{
		userRepo:       userRepo,
		contextTimeout: timeoutContext,
		db:             db,
	}
}

func (u *UserUsecase) FindAll(ctx *gin.Context, params models.FindAllUserParams) ([]*models.User, *types.Error) {
	result, err := u.userRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".UserUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *UserUsecase) Find(ctx *gin.Context, id string) (*models.User, *types.Error) {
	result, err := u.userRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".UserUsecase->Find()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *UserUsecase) Count(ctx *gin.Context, params models.FindAllUserParams) (int, *types.Error) {
	result, err := u.userRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".UserUsecase->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}

func (u *UserUsecase) Create(ctx *gin.Context, obj models.User) (*models.User, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".UserUsecase->Create()" + err.Path
		return nil, err
	}

	data := models.User{}
	data.ID = uuid.New().String()
	data.Name = obj.Name
	data.Email = obj.Email
	data.Username = obj.Username
	data.Password = obj.Password
	data.BusinessID = obj.BusinessID
	data.StatusID = models.DEFAULT_STATUS_CODE

	result, err := u.userRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".UserUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *UserUsecase) Update(ctx *gin.Context, id string, obj models.User) (*models.User, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".UserUsecase->Update()" + err.Path
		return nil, err
	}

	data, err := u.userRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".UserUsecase->Update()" + err.Path
		return nil, err
	}

	data.Name = obj.Name
	data.Email = obj.Email
	data.Username = obj.Username
	data.BusinessID = obj.BusinessID
	data.StatusID = obj.StatusID

	result, err := u.userRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".UserUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *UserUsecase) UpdatePassword(ctx *gin.Context, id string, newPassword string) (*models.User, *types.Error) {
	data, err := u.userRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".UserUsecase->UpdatePassword()" + err.Path
		return nil, err
	}

	hash := md5.New()
	io.WriteString(hash, newPassword)
	data.Password = fmt.Sprintf("%x", hash.Sum(nil))

	result, err := u.userRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".UserUsecase->UpdatePassword()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *UserUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.User, *types.Error) {
	if newStatusID != models.STATUS_ACTIVE && newStatusID != models.STATUS_INACTIVE {
		return nil, &types.Error{
			Path:       ".UserUsecase->UpdateStatus()",
			Message:    "StatusID is not valid",
			Error:      fmt.Errorf("StatusID is not valid"),
			StatusCode: http.StatusBadRequest,
		}
	}

	result, err := u.userRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".UserUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *UserUsecase) Login(ctx *gin.Context, creds models.UserLogin) (*models.UserLogin, *types.Error) {
	err := helpers.ValidateStruct(creds)
	if err != nil {
		err.Path = ".UserUsecase->Login()" + err.Path
		return nil, err
	}

	var userParams models.FindAllUserParams
	userParams.Username = creds.Username
	userParams.Password = creds.Password
	userParams.FindAllParams.StatusID = `status_id = 1`
	users, err := u.FindAll(ctx, userParams)
	if err != nil {
		err.Path = ".UserUsecase->Login()" + err.Path
		return nil, err
	}

	if len(users) == 0 {
		return nil, &types.Error{
			StatusCode: http.StatusUnauthorized,
			Message:    "Username / Password is incorrect",
			Error:      data.ErrNotFound,
			Path:       ".UserUsecase->Login()",
		}
	}

	user := users[0]

	credentials := library.Credential{ID: user.ID, Username: user.Username, Name: user.Name, BusinessID: user.BusinessID, Type: "WebAdmin"}

	token, errorJwtSign := library.JwtSignString(credentials)
	if errorJwtSign != nil {
		return nil, &types.Error{
			Error:      errorJwtSign,
			Message:    "Error JWT Sign String",
			Path:       ".UserUsecase->Login()",
			StatusCode: http.StatusInternalServerError,
		}
	}

	creds.Name = user.Name
	creds.Token = token
	creds.Password = ""

	return &creds, nil
}
