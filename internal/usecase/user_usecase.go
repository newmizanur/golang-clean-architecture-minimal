package usecase

import (
	"context"
	"golang-clean-architecture/internal/apperror"
	"golang-clean-architecture/internal/auth"
	"golang-clean-architecture/internal/dto"
	"golang-clean-architecture/internal/dto/converter"
	m "golang-clean-architecture/internal/persistence/model"
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
	UserRepository UserRepositoryPort
	JwtSecret      string
	JwtTTL         time.Duration
}

func NewUserUseCase(db *bun.DB, logger *logrus.Logger, validate *validator.Validate,
	userRepository UserRepositoryPort, jwtSecret string, jwtTTL time.Duration) *UserUseCase {
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
		c.Log.WithError(err).Error("Failed to start transaction")
		return nil, apperror.UserErrors.FailedToCreate
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("Invalid request body")
		return nil, apperror.UserErrors.InvalidRequest
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.WithError(err).Error("Failed to generate bcrypt hash")
		return nil, apperror.UserErrors.FailedToCreate
	}

	now := time.Now().UnixMilli()
	user := &m.User{
		ID:        uuid.NewString(),
		Password:  string(password),
		Name:      request.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.UserRepository.Create(ctx, tx, user); err != nil {
		c.Log.WithError(err).Error("Failed to create user in database")
		return nil, apperror.UserErrors.FailedToCreate
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, apperror.UserErrors.FailedToCreate
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *dto.LoginUserRequest) (*dto.UserResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("Invalid request body")
		return nil, apperror.UserErrors.InvalidRequest
	}

	// Read-only operation, no transaction needed
	user, err := c.UserRepository.FindById(ctx, nil, request.ID)
	if err != nil {
		c.Log.WithError(err).WithField("user_id", request.ID).Warn("Failed to find user by id")
		return nil, apperror.UserErrors.Unauthorized
	}
	if user == nil {
		c.Log.WithField("user_id", request.ID).Warn("User not found")
		return nil, apperror.UserErrors.Unauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.WithError(err).Warn("Password mismatch")
		return nil, apperror.UserErrors.Unauthorized
	}

	jwtToken, err := auth.GenerateToken(user.ID, c.JwtSecret, c.JwtTTL)
	if err != nil {
		c.Log.WithError(err).Error("Failed to generate JWT token")
		return nil, apperror.UserErrors.FailedToLogin
	}

	return converter.UserToTokenResponse(jwtToken), nil
}

func (c *UserUseCase) Current(ctx context.Context, request *dto.GetUserRequest) (*dto.UserResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("Invalid request body")
		return nil, apperror.UserErrors.InvalidRequest
	}

	// Read-only operation, no transaction needed
	user, err := c.UserRepository.FindById(ctx, nil, request.ID)
	if err != nil {
		c.Log.WithError(err).WithField("user_id", request.ID).Error("Failed to find user by id")
		return nil, apperror.UserErrors.NotFound
	}
	if user == nil {
		c.Log.WithField("user_id", request.ID).Warn("User not found")
		return nil, apperror.UserErrors.NotFound
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Logout(ctx context.Context, request *dto.LogoutUserRequest) (bool, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("Invalid request body")
		return false, apperror.UserErrors.InvalidRequest
	}

	// JWT logout is client-side (token removal)
	// No server-side state to update
	c.Log.WithField("user_id", request.ID).Info("User logged out")
	return true, nil
}

func (c *UserUseCase) Update(ctx context.Context, request *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("Failed to start transaction")
		return nil, apperror.UserErrors.FailedToUpdate
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("Invalid request body")
		return nil, apperror.UserErrors.InvalidRequest
	}

	user, err := c.UserRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		c.Log.WithError(err).WithField("user_id", request.ID).Error("Failed to find user by id")
		return nil, apperror.UserErrors.NotFound
	}
	if user == nil {
		c.Log.WithField("user_id", request.ID).Warn("User not found")
		return nil, apperror.UserErrors.NotFound
	}

	// Update fields - OmitZero in repository will handle partial updates
	user.Name = request.Name

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.WithError(err).Error("Failed to generate bcrypt hash")
			return nil, apperror.UserErrors.FailedToUpdate
		}
		user.Password = string(password)
	}

	user.UpdatedAt = time.Now().UnixMilli()
	if err := c.UserRepository.Update(ctx, tx, user); err != nil {
		c.Log.WithError(err).Error("Failed to update user")
		return nil, apperror.UserErrors.FailedToUpdate
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, apperror.UserErrors.FailedToUpdate
	}

	return converter.UserToResponse(user), nil
}
