package usecase

import (
	"context"
	"golang-clean-architecture/internal/apperror"
	"golang-clean-architecture/internal/auth"
	"golang-clean-architecture/internal/dto"
	"golang-clean-architecture/internal/dto/converter"
	dbmodel "golang-clean-architecture/internal/persistence/model"
	"golang-clean-architecture/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	DB             *bun.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	JwtSecret      string
	JwtTTL         time.Duration
}

func NewUserUseCase(db *bun.DB, logger *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository, jwtSecret string, jwtTTL time.Duration) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
		JwtSecret:      jwtSecret,
		JwtTTL:         jwtTTL,
	}
}

func (c *UserUseCase) Create(ctx context.Context, request *dto.RegisterUserRequest) (*dto.UserResponse, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.Warnf("Failed to start transaction : %+v", err)
		return nil, apperror.UserErrors.FailedToCreate
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, apperror.UserErrors.InvalidRequest
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
		return nil, apperror.UserErrors.FailedToCreate
	}

	now := time.Now().UnixMilli()
	user := &dbmodel.Users{
		ID:        uuid.NewString(),
		Password:  string(password),
		Name:      request.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.UserRepository.Create(ctx, tx, user); err != nil {
		c.Log.Warnf("Failed create user to database : %+v", err)
		return nil, apperror.UserErrors.FailedToCreate
	}

	if err := tx.Commit(); err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, apperror.UserErrors.FailedToCreate
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *dto.LoginUserRequest) (*dto.UserResponse, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.Warnf("Failed to start transaction : %+v", err)
		return nil, apperror.UserErrors.FailedToLogin
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body  : %+v", err)
		return nil, apperror.UserErrors.InvalidRequest
	}

	user, err := c.UserRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, apperror.UserErrors.Unauthorized
	}
	if user == nil {
		c.Log.Warnf("User not found : %+v", request.ID)
		return nil, apperror.UserErrors.Unauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Failed to compare user password with bcrype hash : %+v", err)
		return nil, apperror.UserErrors.Unauthorized
	}

	jwtToken, err := auth.GenerateToken(user.ID, c.JwtSecret, c.JwtTTL)
	if err != nil {
		c.Log.Warnf("Failed to generate JWT token : %+v", err)
		return nil, apperror.UserErrors.FailedToLogin
	}
	user.UpdatedAt = time.Now().UnixMilli()
	if err := c.UserRepository.Update(ctx, tx, user); err != nil {
		c.Log.Warnf("Failed save user : %+v", err)
		return nil, apperror.UserErrors.FailedToLogin
	}

	if err := tx.Commit(); err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, apperror.UserErrors.FailedToLogin
	}

	return converter.UserToTokenResponse(jwtToken), nil
}

func (c *UserUseCase) Current(ctx context.Context, request *dto.GetUserRequest) (*dto.UserResponse, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.Warnf("Failed to start transaction : %+v", err)
		return nil, apperror.UserErrors.FailedToGet
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, apperror.UserErrors.InvalidRequest
	}

	user, err := c.UserRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, apperror.UserErrors.NotFound
	}
	if user == nil {
		c.Log.Warnf("User not found : %+v", request.ID)
		return nil, apperror.UserErrors.NotFound
	}

	if err := tx.Commit(); err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, apperror.UserErrors.FailedToGet
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Logout(ctx context.Context, request *dto.LogoutUserRequest) (bool, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.Warnf("Failed to start transaction : %+v", err)
		return false, apperror.UserErrors.FailedToLogout
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, apperror.UserErrors.InvalidRequest
	}

	user, err := c.UserRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return false, apperror.UserErrors.NotFound
	}
	if user == nil {
		c.Log.Warnf("User not found : %+v", request.ID)
		return false, apperror.UserErrors.NotFound
	}

	user.UpdatedAt = time.Now().UnixMilli()
	if err := c.UserRepository.Update(ctx, tx, user); err != nil {
		c.Log.Warnf("Failed save user : %+v", err)
		return false, apperror.UserErrors.FailedToLogout
	}

	if err := tx.Commit(); err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, apperror.UserErrors.FailedToLogout
	}

	return true, nil
}

func (c *UserUseCase) Update(ctx context.Context, request *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.Warnf("Failed to start transaction : %+v", err)
		return nil, apperror.UserErrors.FailedToUpdate
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, apperror.UserErrors.InvalidRequest
	}

	user, err := c.UserRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, apperror.UserErrors.NotFound
	}
	if user == nil {
		c.Log.Warnf("User not found : %+v", request.ID)
		return nil, apperror.UserErrors.NotFound
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
			return nil, apperror.UserErrors.FailedToUpdate
		}
		user.Password = string(password)
	}

	user.UpdatedAt = time.Now().UnixMilli()
	if err := c.UserRepository.Update(ctx, tx, user); err != nil {
		c.Log.Warnf("Failed save user : %+v", err)
		return nil, apperror.UserErrors.FailedToUpdate
	}

	if err := tx.Commit(); err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, apperror.UserErrors.FailedToUpdate
	}

	return converter.UserToResponse(user), nil
}
