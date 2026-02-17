package usecase

import (
	"context"
	"golang-clean-architecture/internal/apperror"
	"golang-clean-architecture/internal/dto"
	"golang-clean-architecture/internal/dto/converter"
	m "golang-clean-architecture/internal/persistence/model"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type ContactUseCase struct {
	DB                *bun.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	ContactRepository ContactRepositoryPort
}

func NewContactUseCase(db *bun.DB, logger *logrus.Logger, validate *validator.Validate,
	contactRepository ContactRepositoryPort) *ContactUseCase {
	return &ContactUseCase{
		DB:                db,
		Log:               logger,
		Validate:          validate,
		ContactRepository: contactRepository,
	}
}

func (c *ContactUseCase) Create(ctx context.Context, request *dto.CreateContactRequest) (*dto.ContactResponse, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("error starting transaction")
		return nil, apperror.ContactErrors.FailedToCreate
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, apperror.ContactErrors.InvalidRequest
	}

	now := time.Now().UnixMilli()
	contact := &m.Contact{
		ID:        uuid.NewString(),
		FirstName: request.FirstName,
		LastName:  stringPtrOrNil(request.LastName),
		Email:     stringPtrOrNil(request.Email),
		Phone:     stringPtrOrNil(request.Phone),
		UserID:    request.UserId,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.ContactRepository.Create(ctx, tx, contact); err != nil {
		c.Log.WithError(err).Error("error creating contact")
		return nil, apperror.ContactErrors.FailedToCreate
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("error creating contact")
		return nil, apperror.ContactErrors.FailedToCreate
	}

	return converter.ContactToResponse(contact), nil
}

func (c *ContactUseCase) Update(ctx context.Context, request *dto.UpdateContactRequest) (*dto.ContactResponse, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("error starting transaction")
		return nil, apperror.ContactErrors.FailedToUpdate
	}
	defer tx.Rollback()

	contact, err := c.ContactRepository.FindByIdAndUserId(ctx, tx, request.ID, request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, apperror.ContactErrors.NotFound
	}
	if contact == nil {
		c.Log.WithField("contact_id", request.ID).Warn("contact not found")
		return nil, apperror.ContactErrors.NotFound
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, apperror.ContactErrors.InvalidRequest
	}

	contact.FirstName = request.FirstName
	contact.LastName = stringPtrOrNil(request.LastName)
	contact.Email = stringPtrOrNil(request.Email)
	contact.Phone = stringPtrOrNil(request.Phone)
	contact.UpdatedAt = time.Now().UnixMilli()

	if err := c.ContactRepository.Update(ctx, tx, contact); err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, apperror.ContactErrors.FailedToUpdate
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, apperror.ContactErrors.FailedToUpdate
	}

	return converter.ContactToResponse(contact), nil
}

func (c *ContactUseCase) Get(ctx context.Context, request *dto.GetContactRequest) (*dto.ContactResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, apperror.ContactErrors.InvalidRequest
	}

	// Read-only operation, no transaction needed
	contact, err := c.ContactRepository.FindByIdAndUserId(ctx, nil, request.ID, request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, apperror.ContactErrors.NotFound
	}
	if contact == nil {
		c.Log.WithField("contact_id", request.ID).Warn("contact not found")
		return nil, apperror.ContactErrors.NotFound
	}

	return converter.ContactToResponse(contact), nil
}

func (c *ContactUseCase) Delete(ctx context.Context, request *dto.DeleteContactRequest) error {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("error starting transaction")
		return apperror.ContactErrors.FailedToDelete
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return apperror.ContactErrors.InvalidRequest
	}

	contact, err := c.ContactRepository.FindByIdAndUserId(ctx, tx, request.ID, request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return apperror.ContactErrors.NotFound
	}
	if contact == nil {
		c.Log.WithField("contact_id", request.ID).Warn("contact not found")
		return apperror.ContactErrors.NotFound
	}

	if err := c.ContactRepository.Delete(ctx, tx, contact); err != nil {
		c.Log.WithError(err).Error("error deleting contact")
		return apperror.ContactErrors.FailedToDelete
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("error deleting contact")
		return apperror.ContactErrors.FailedToDelete
	}

	return nil
}

func (c *ContactUseCase) Search(ctx context.Context, request *dto.SearchContactRequest) ([]dto.ContactResponse, int64, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, apperror.ContactErrors.InvalidRequest
	}

	// Read-only operation, no transaction needed
	contacts, total, err := c.ContactRepository.Search(ctx, nil, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting contacts")
		return nil, 0, apperror.ContactErrors.FailedToSearch
	}

	responses := make([]dto.ContactResponse, len(contacts))
	for i, contact := range contacts {
		responses[i] = *converter.ContactToResponse(&contact)
	}

	return responses, total, nil
}
