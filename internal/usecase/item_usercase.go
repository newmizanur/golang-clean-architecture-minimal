package usecase

import (
	"context"
	"time"

	"golang-clean-architecture/ent"
	"golang-clean-architecture/internal/apperror"
	"golang-clean-architecture/internal/dto"
	"golang-clean-architecture/internal/dto/converter"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ItemUseCase struct {
	Client         *ent.Client
	Log            *logrus.Logger
	Validate       *validator.Validate
	ItemRepository ItemRepositoryPort
}

func NewItemUseCase(client *ent.Client, log *logrus.Logger, validate *validator.Validate, repo ItemRepositoryPort) *ItemUseCase {
	return &ItemUseCase{
		Client:         client,
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

	tx, err := c.Client.Tx(ctx)
	if err != nil {
		c.Log.WithError(err).Error("error on starting transaction at item usecase")
		return nil, apperror.ItemErrors.FailedToCreateTransaction
	}
	defer tx.Rollback()

	now := time.Now()
	item := &ent.Item{
		Name:      request.Name,
		Sku:       request.SKU,
		Currency:  request.Currency,
		Stock:     request.Stock,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, err := c.ItemRepository.Create(ctx, tx.Client(), item)
	if err != nil {
		c.Log.WithError(err).Error("error on creating item")
		return nil, apperror.ItemErrors.FailedToCreateItem
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("error creating item")
		return nil, apperror.ItemErrors.FailedToCreateItem
	}

	return converter.ItemToResponse(created), nil
}

func (c *ItemUseCase) Search(ctx context.Context, request *dto.SearchItemRequest) ([]dto.CreateItemResponse, int64, error) {
	items, total, err := c.ItemRepository.Search(ctx, c.Client, request)
	if err != nil {
		return nil, 0, err
	}

	response := make([]dto.CreateItemResponse, len(items))
	for i, item := range items {
		response[i] = *converter.ItemToResponse(item)
	}

	return response, total, nil
}

func (c *ItemUseCase) Get(ctx context.Context, request *dto.GetItemRequest) (*dto.CreateItemResponse, error) {
	item, err := c.ItemRepository.FindById(ctx, c.Client, int(request.ID))
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

	tx, err := c.Client.Tx(ctx)
	if err != nil {
		c.Log.WithError(err).Error("error starting transaction on update item")
		return nil, apperror.ItemErrors.FailedToUpdate
	}
	defer tx.Rollback()

	item, err := c.ItemRepository.FindById(ctx, tx.Client(), int(request.ID))
	if err != nil {
		c.Log.WithError(err).Error("error getting item")
		return nil, apperror.ItemErrors.FailedToUpdate
	}
	if item == nil {
		return nil, apperror.ItemErrors.NotFound
	}

	item.Name = request.Name
	item.Sku = request.SKU
	item.Currency = request.Currency
	if request.Stock >= 0 {
		item.Stock = request.Stock
	}
	item.UpdatedAt = time.Now()

	if err := c.ItemRepository.Update(ctx, tx.Client(), item); err != nil {
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
	tx, err := c.Client.Tx(ctx)
	if err != nil {
		c.Log.WithError(err).Error("error starting transaction on delete item")
		return apperror.ItemErrors.FailedToDelete
	}
	defer tx.Rollback()

	item, err := c.ItemRepository.FindById(ctx, tx.Client(), int(request.ID))
	if err != nil {
		c.Log.WithError(err).Error("error getting item")
		return apperror.ItemErrors.FailedToDelete
	}
	if item == nil {
		return apperror.ItemErrors.NotFound
	}

	if err := c.ItemRepository.Delete(ctx, tx.Client(), item.ID); err != nil {
		c.Log.WithError(err).Error("error deleting item")
		return apperror.ItemErrors.FailedToDelete
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("error commit transaction on delete item")
		return apperror.ItemErrors.FailedToDelete
	}

	return nil
}
