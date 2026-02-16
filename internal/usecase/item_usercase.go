package usecase

import (
	"context"
	"time"

	"golang-clean-architecture/internal/apperror"
	"golang-clean-architecture/internal/dto"
	"golang-clean-architecture/internal/dto/converter"
	m "golang-clean-architecture/internal/persistence/model"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type ItemUseCase struct {
	DB             *bun.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	ItemRepository ItemRepositoryPort
}

func NewItemUseCase(db *bun.DB, log *logrus.Logger, validate *validator.Validate, repo ItemRepositoryPort) *ItemUseCase {
	return &ItemUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		ItemRepository: repo,
	}
}

func (c *ItemUseCase) Create(ctx context.Context, request *dto.CreateItemRequest) (*dto.CreateItemResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("invalid request body")
		return nil, apperror.ItemErrors.InvalidRequest
	}

	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("error on starting transaction at item usecase")
		return nil, apperror.ItemErrors.FailedToCreateTransaction
	}
	defer tx.Rollback()

	now := time.Now()
	item := m.Items{
		Name:      request.Name,
		Sku:       request.SKU,
		Currency:  request.Currency,
		Stock:     request.Stock,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	id, err := c.ItemRepository.Create(ctx, tx, &item)
	if err != nil {
		c.Log.WithError(err).Error("error on creating item")
		return nil, apperror.ItemErrors.FailedToCreateItem
	}
	item.ID = id

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("error creating item")
		return nil, apperror.ItemErrors.FailedToCreateItem
	}

	return converter.ItemToResponse(&item), nil
}

func (c *ItemUseCase) Search(ctx context.Context, request *dto.SearchItemRequest) ([]dto.CreateItemResponse, int64, error) {
	var response []dto.CreateItemResponse
	// Read-only operation, no transaction needed
	items, total, err := c.ItemRepository.Search(ctx, nil, request)
	if err != nil {
		return response, 0, err
	}

	response = make([]dto.CreateItemResponse, len(items))
	for i, item := range items {
		response[i] = *converter.ItemToResponse(&item)
	}

	return response, total, nil
}

func (c *ItemUseCase) Get(ctx context.Context, request *dto.GetItemRequest) (*dto.CreateItemResponse, error) {
	// Read-only operation, no transaction needed
	item, err := c.ItemRepository.FindById(ctx, nil, request.ID)
	if err != nil {
		c.Log.WithError(err).Error("error getting item")
		return nil, apperror.ItemErrors.FailedToGet
	}
	if item == nil {
		return nil, apperror.ItemErrors.NotFound
	}

	return converter.ItemToResponse(item), nil
}

func (c *ItemUseCase) Update(ctx context.Context, request *dto.UpdateItemRequest) (*dto.CreateItemResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("invalid request body")
		return nil, apperror.ItemErrors.InvalidRequest
	}

	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("error starting transaction on update item")
		return nil, apperror.ItemErrors.FailedToUpdate
	}
	defer tx.Rollback()

	item, err := c.ItemRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		c.Log.WithError(err).Error("error getting item")
		return nil, apperror.ItemErrors.FailedToUpdate
	}
	if item == nil {
		return nil, apperror.ItemErrors.NotFound
	}

	// Update fields - OmitZero in repository will handle partial updates
	item.Name = request.Name
	item.Sku = request.SKU
	item.Currency = request.Currency
	if request.Stock >= 0 {
		item.Stock = request.Stock
	}

	now := time.Now()
	item.UpdatedAt = &now

	if err := c.ItemRepository.Update(ctx, tx, item); err != nil {
		c.Log.WithError(err).Error("error updating item")
		return nil, apperror.ItemErrors.FailedToUpdate
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("error commit transaction on update item")
		return nil, apperror.ItemErrors.FailedToUpdate
	}

	return converter.ItemToResponse(item), nil
}

func (c *ItemUseCase) Delete(ctx context.Context, request *dto.DeleteItemRequest) error {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("error starting transaction on delete item")
		return apperror.ItemErrors.FailedToDelete
	}
	defer tx.Rollback()

	item, err := c.ItemRepository.FindById(ctx, tx, request.ID)
	if err != nil {
		c.Log.WithError(err).Error("error getting item")
		return apperror.ItemErrors.FailedToDelete
	}
	if item == nil {
		return apperror.ItemErrors.NotFound
	}

	if err := c.ItemRepository.Delete(ctx, tx, item); err != nil {
		c.Log.WithError(err).Error("error deleting item")
		return apperror.ItemErrors.FailedToDelete
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("error commit transaction on delete item")
		return apperror.ItemErrors.FailedToDelete
	}

	return nil
}
