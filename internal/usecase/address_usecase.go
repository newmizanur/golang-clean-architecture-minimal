package usecase

import (
	"context"
	"golang-clean-architecture/ent"
	"golang-clean-architecture/internal/apperror"
	"golang-clean-architecture/internal/dto"
	"golang-clean-architecture/internal/dto/converter"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AddressUseCase struct {
	Client            *ent.Client
	Log               *logrus.Logger
	Validate          *validator.Validate
	AddressRepository AddressRepositoryPort
	ContactRepository ContactRepositoryPort
}

func NewAddressUseCase(client *ent.Client, logger *logrus.Logger, validate *validator.Validate,
	contactRepository ContactRepositoryPort, addressRepository AddressRepositoryPort) *AddressUseCase {
	return &AddressUseCase{
		Client:            client,
		Log:               logger,
		Validate:          validate,
		ContactRepository: contactRepository,
		AddressRepository: addressRepository,
	}
}

func (c *AddressUseCase) Create(ctx context.Context, request *dto.CreateAddressRequest) (*dto.AddressResponse, error) {
	tx, err := c.Client.Tx(ctx)
	if err != nil {
		c.Log.WithError(err).Error("failed to start transaction")
		return nil, apperror.AddressErrors.FailedToCreate
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, apperror.AddressErrors.InvalidRequest
	}

	contact, err := c.ContactRepository.FindByIdAndUserId(ctx, tx.Client(), request.ContactId, request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return nil, apperror.AddressErrors.FailedToCreate
	}
	if contact == nil {
		c.Log.WithField("contact_id", request.ContactId).Warn("contact not found")
		return nil, apperror.ContactErrors.NotFound
	}

	now := time.Now().UnixMilli()
	address := &ent.Address{
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

	if err := c.AddressRepository.Create(ctx, tx.Client(), address); err != nil {
		c.Log.WithError(err).Error("failed to create address")
		return nil, apperror.AddressErrors.FailedToCreate
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, apperror.AddressErrors.FailedToCreate
	}

	return converter.AddressToResponse(address), nil
}

func (c *AddressUseCase) Update(ctx context.Context, request *dto.UpdateAddressRequest) (*dto.AddressResponse, error) {
	tx, err := c.Client.Tx(ctx)
	if err != nil {
		c.Log.WithError(err).Error("failed to start transaction")
		return nil, apperror.AddressErrors.FailedToUpdate
	}
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, apperror.AddressErrors.InvalidRequest
	}

	contact, err := c.ContactRepository.FindByIdAndUserId(ctx, tx.Client(), request.ContactId, request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return nil, apperror.AddressErrors.FailedToUpdate
	}
	if contact == nil {
		c.Log.WithField("contact_id", request.ContactId).Warn("contact not found")
		return nil, apperror.ContactErrors.NotFound
	}

	address, err := c.AddressRepository.FindByIdAndContactId(ctx, tx.Client(), request.ID, contact.ID)
	if err != nil {
		c.Log.WithError(err).Error("failed to find address")
		return nil, apperror.AddressErrors.FailedToUpdate
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

	if err := c.AddressRepository.Update(ctx, tx.Client(), address); err != nil {
		c.Log.WithError(err).Error("failed to update address")
		return nil, apperror.AddressErrors.FailedToUpdate
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, apperror.AddressErrors.FailedToUpdate
	}

	return converter.AddressToResponse(address), nil
}

func (c *AddressUseCase) Get(ctx context.Context, request *dto.GetAddressRequest) (*dto.AddressResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, apperror.AddressErrors.InvalidRequest
	}

	contact, err := c.ContactRepository.FindByIdAndUserId(ctx, c.Client, request.ContactId, request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return nil, apperror.AddressErrors.FailedToGet
	}
	if contact == nil {
		c.Log.WithField("contact_id", request.ContactId).Warn("contact not found")
		return nil, apperror.ContactErrors.NotFound
	}

	address, err := c.AddressRepository.FindByIdAndContactId(ctx, c.Client, request.ID, contact.ID)
	if err != nil {
		c.Log.WithError(err).Error("failed to find address")
		return nil, apperror.AddressErrors.FailedToGet
	}
	if address == nil {
		c.Log.WithField("address_id", request.ID).Warn("address not found")
		return nil, apperror.AddressErrors.NotFound
	}

	return converter.AddressToResponse(address), nil
}

func (c *AddressUseCase) Delete(ctx context.Context, request *dto.DeleteAddressRequest) error {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return apperror.AddressErrors.InvalidRequest
	}

	tx, err := c.Client.Tx(ctx)
	if err != nil {
		c.Log.WithError(err).Error("failed to start transaction")
		return apperror.AddressErrors.FailedToDelete
	}
	defer tx.Rollback()

	contact, err := c.ContactRepository.FindByIdAndUserId(ctx, tx.Client(), request.ContactId, request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return apperror.AddressErrors.FailedToDelete
	}
	if contact == nil {
		c.Log.WithField("contact_id", request.ContactId).Warn("contact not found")
		return apperror.ContactErrors.NotFound
	}

	address, err := c.AddressRepository.FindByIdAndContactId(ctx, tx.Client(), request.ID, contact.ID)
	if err != nil {
		c.Log.WithError(err).Error("failed to find address")
		return apperror.AddressErrors.FailedToDelete
	}
	if address == nil {
		c.Log.WithField("address_id", request.ID).Warn("address not found")
		return apperror.AddressErrors.NotFound
	}

	if err := c.AddressRepository.Delete(ctx, tx.Client(), address.ID); err != nil {
		c.Log.WithError(err).Error("failed to delete address")
		return apperror.AddressErrors.FailedToDelete
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return apperror.AddressErrors.FailedToDelete
	}

	return nil
}

func (c *AddressUseCase) List(ctx context.Context, request *dto.ListAddressRequest) ([]dto.AddressResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, apperror.AddressErrors.InvalidRequest
	}

	contact, err := c.ContactRepository.FindByIdAndUserId(ctx, c.Client, request.ContactId, request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return nil, apperror.AddressErrors.FailedToList
	}
	if contact == nil {
		c.Log.WithField("contact_id", request.ContactId).Warn("contact not found")
		return nil, apperror.ContactErrors.NotFound
	}

	addresses, err := c.AddressRepository.FindAllByContactId(ctx, c.Client, contact.ID)
	if err != nil {
		c.Log.WithError(err).Error("failed to find addresses")
		return nil, apperror.AddressErrors.FailedToList
	}

	responses := make([]dto.AddressResponse, len(addresses))
	for i, address := range addresses {
		responses[i] = *converter.AddressToResponse(address)
	}

	return responses, nil
}
