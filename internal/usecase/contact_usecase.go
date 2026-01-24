package usecase

import (
	"context"
	"database/sql"
	"golang-clean-architecture/internal/apperror"
	dbmodel "golang-clean-architecture/internal/entity/db/model"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ContactUseCase struct {
	DB                *sql.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	ContactRepository *repository.ContactRepository
}

func NewContactUseCase(db *sql.DB, logger *logrus.Logger, validate *validator.Validate,
	contactRepository *repository.ContactRepository) *ContactUseCase {
	return &ContactUseCase{
		DB:                db,
		Log:               logger,
		Validate:          validate,
		ContactRepository: contactRepository,
	}
}

func (c *ContactUseCase) Create(ctx context.Context, request *model.CreateContactRequest) (*model.ContactResponse, error) {
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
	contact := &dbmodel.Contacts{
		ID:        uuid.New().String(),
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

func (c *ContactUseCase) Update(ctx context.Context, request *model.UpdateContactRequest) (*model.ContactResponse, error) {
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

func (c *ContactUseCase) Get(ctx context.Context, request *model.GetContactRequest) (*model.ContactResponse, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("error starting transaction")
		return nil, apperror.ContactErrors.FailedToGet
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, apperror.ContactErrors.InvalidRequest
	}

	contact, err := c.ContactRepository.FindByIdAndUserId(ctx, tx, request.ID, request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, apperror.ContactErrors.NotFound
	}
	if contact == nil {
		c.Log.WithField("contact_id", request.ID).Warn("contact not found")
		return nil, apperror.ContactErrors.NotFound
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, apperror.ContactErrors.FailedToGet
	}

	return converter.ContactToResponse(contact), nil
}

func (c *ContactUseCase) Delete(ctx context.Context, request *model.DeleteContactRequest) error {
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

func (c *ContactUseCase) Search(ctx context.Context, request *model.SearchContactRequest) ([]model.ContactResponse, int64, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("error starting transaction")
		return nil, 0, apperror.ContactErrors.FailedToSearch
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, apperror.ContactErrors.InvalidRequest
	}

	contacts, total, err := c.ContactRepository.Search(ctx, tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting contacts")
		return nil, 0, apperror.ContactErrors.FailedToSearch
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("error getting contacts")
		return nil, 0, apperror.ContactErrors.FailedToSearch
	}

	responses := make([]model.ContactResponse, len(contacts))
	for i, contact := range contacts {
		responses[i] = *converter.ContactToResponse(&contact)
	}

	return responses, total, nil
}
