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

type AddressUseCase struct {
	DB                *sql.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	AddressRepository *repository.AddressRepository
	ContactRepository *repository.ContactRepository
}

func NewAddressUseCase(db *sql.DB, logger *logrus.Logger, validate *validator.Validate,
	contactRepository *repository.ContactRepository, addressRepository *repository.AddressRepository) *AddressUseCase {
	return &AddressUseCase{
		DB:                db,
		Log:               logger,
		Validate:          validate,
		ContactRepository: contactRepository,
		AddressRepository: addressRepository,
	}
}

func (c *AddressUseCase) Create(ctx context.Context, request *model.CreateAddressRequest) (*model.AddressResponse, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("failed to start transaction")
		return nil, apperror.AddressErrors.FailedToCreate
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, apperror.AddressErrors.InvalidRequest
	}

	contact, err := c.ContactRepository.FindByIdAndUserId(ctx, tx, request.ContactId, request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return nil, apperror.AddressErrors.NotFound
	}
	if contact == nil {
		c.Log.WithField("contact_id", request.ContactId).Warn("contact not found")
		return nil, apperror.AddressErrors.NotFound
	}

	now := time.Now().UnixMilli()
	address := &dbmodel.Addresses{
		ID:         uuid.NewString(),
		ContactID:  contact.ID,
		Street:     stringPtrOrNil(request.Street),
		City:       stringPtrOrNil(request.City),
		Province:   stringPtrOrNil(request.Province),
		PostalCode: stringPtrOrNil(request.PostalCode),
		Country:    stringPtrOrNil(request.Country),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := c.AddressRepository.Create(ctx, tx, address); err != nil {
		c.Log.WithError(err).Error("failed to create address")
		return nil, apperror.AddressErrors.FailedToCreate
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, apperror.AddressErrors.FailedToCreate
	}

	return converter.AddressToResponse(address), nil
}

func (c *AddressUseCase) Update(ctx context.Context, request *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("failed to start transaction")
		return nil, apperror.AddressErrors.FailedToUpdate
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, apperror.AddressErrors.InvalidRequest
	}

	contact, err := c.ContactRepository.FindByIdAndUserId(ctx, tx, request.ContactId, request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return nil, apperror.AddressErrors.NotFound
	}
	if contact == nil {
		c.Log.WithField("contact_id", request.ContactId).Warn("contact not found")
		return nil, apperror.AddressErrors.NotFound
	}

	address, err := c.AddressRepository.FindByIdAndContactId(ctx, tx, request.ID, contact.ID)
	if err != nil {
		c.Log.WithError(err).Error("failed to find address")
		return nil, apperror.AddressErrors.NotFound
	}
	if address == nil {
		c.Log.WithField("address_id", request.ID).Warn("address not found")
		return nil, apperror.AddressErrors.NotFound
	}

	address.Street = stringPtrOrNil(request.Street)
	address.City = stringPtrOrNil(request.City)
	address.Province = stringPtrOrNil(request.Province)
	address.PostalCode = stringPtrOrNil(request.PostalCode)
	address.Country = stringPtrOrNil(request.Country)
	address.UpdatedAt = time.Now().UnixMilli()

	if err := c.AddressRepository.Update(ctx, tx, address); err != nil {
		c.Log.WithError(err).Error("failed to update address")
		return nil, apperror.AddressErrors.FailedToUpdate
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, apperror.AddressErrors.FailedToUpdate
	}

	return converter.AddressToResponse(address), nil
}

func (c *AddressUseCase) Get(ctx context.Context, request *model.GetAddressRequest) (*model.AddressResponse, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("failed to start transaction")
		return nil, apperror.AddressErrors.FailedToGet
	}
	defer tx.Rollback()

	contact, err := c.ContactRepository.FindByIdAndUserId(ctx, tx, request.ContactId, request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return nil, apperror.AddressErrors.NotFound
	}
	if contact == nil {
		c.Log.WithField("contact_id", request.ContactId).Warn("contact not found")
		return nil, apperror.AddressErrors.NotFound
	}

	address, err := c.AddressRepository.FindByIdAndContactId(ctx, tx, request.ID, contact.ID)
	if err != nil {
		c.Log.WithError(err).Error("failed to find address")
		return nil, apperror.AddressErrors.NotFound
	}
	if address == nil {
		c.Log.WithField("address_id", request.ID).Warn("address not found")
		return nil, apperror.AddressErrors.NotFound
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, apperror.AddressErrors.FailedToGet
	}

	return converter.AddressToResponse(address), nil
}

func (c *AddressUseCase) Delete(ctx context.Context, request *model.DeleteAddressRequest) error {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("failed to start transaction")
		return apperror.AddressErrors.FailedToDelete
	}
	defer tx.Rollback()

	contact, err := c.ContactRepository.FindByIdAndUserId(ctx, tx, request.ContactId, request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return apperror.AddressErrors.NotFound
	}
	if contact == nil {
		c.Log.WithField("contact_id", request.ContactId).Warn("contact not found")
		return apperror.AddressErrors.NotFound
	}

	address, err := c.AddressRepository.FindByIdAndContactId(ctx, tx, request.ID, contact.ID)
	if err != nil {
		c.Log.WithError(err).Error("failed to find address")
		return apperror.AddressErrors.NotFound
	}
	if address == nil {
		c.Log.WithField("address_id", request.ID).Warn("address not found")
		return apperror.AddressErrors.NotFound
	}

	if err := c.AddressRepository.Delete(ctx, tx, address); err != nil {
		c.Log.WithError(err).Error("failed to delete address")
		return apperror.AddressErrors.FailedToDelete
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return apperror.AddressErrors.FailedToDelete
	}

	return nil
}

func (c *AddressUseCase) List(ctx context.Context, request *model.ListAddressRequest) ([]model.AddressResponse, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("failed to start transaction")
		return nil, apperror.AddressErrors.FailedToList
	}
	defer tx.Rollback()

	contact, err := c.ContactRepository.FindByIdAndUserId(ctx, tx, request.ContactId, request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return nil, apperror.AddressErrors.NotFound
	}
	if contact == nil {
		c.Log.WithField("contact_id", request.ContactId).Warn("contact not found")
		return nil, apperror.AddressErrors.NotFound
	}

	addresses, err := c.AddressRepository.FindAllByContactId(ctx, tx, contact.ID)
	if err != nil {
		c.Log.WithError(err).Error("failed to find addresses")
		return nil, apperror.AddressErrors.FailedToList
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, apperror.AddressErrors.FailedToList
	}

	responses := make([]model.AddressResponse, len(addresses))
	for i, address := range addresses {
		responses[i] = *converter.AddressToResponse(&address)
	}

	return responses, nil
}
