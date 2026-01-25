package usecase

import (
	"context"
	"database/sql"
	"time"

	"golang-clean-architecture/internal/apperror"
	dbmodel "golang-clean-architecture/internal/entity/db/model"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"

	"github.com/sirupsen/logrus"
)

type ItemUseCase struct {
	DB             *sql.DB
	Log            *logrus.Logger
	ItemRepository *repository.ItemRepository
}

func NewItemUseCase(db *sql.DB, log *logrus.Logger, repository *repository.ItemRepository) *ItemUseCase {
	return &ItemUseCase{
		DB:             db,
		Log:            log,
		ItemRepository: repository,
	}
}

func (c *ItemUseCase) Create(ctx context.Context, request *model.CreateItemRequest) (*model.CreateItemResponse, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		c.Log.WithError(err).Error("error on starting transaction at item usecase")
		return nil, apperror.ItemErrors.FailedToCreateTransaction
	}
	defer tx.Rollback()

	now := time.Now()
	item := dbmodel.Items{
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

	it, err := c.ItemRepository.Get(ctx, tx, id)
	if err != nil {
		c.Log.WithError(err).Error("error on get item")
	}

	if err := tx.Commit(); err != nil {
		c.Log.WithError(err).Error("error creating item")
		return nil, apperror.ItemErrors.FailedToCreateItem
	}

	return converter.ItemToResponse(it), nil
}

func (c *ItemUseCase) Search(ctx context.Context, request *model.SearchItemRequest) ([]model.CreateItemResponse, int64, error) {
	var response []model.CreateItemResponse
	items, total, err := c.ItemRepository.Search(ctx, nil, request)
	if err != nil {
		return response, 0, err
	}

	response = make([]model.CreateItemResponse, len(items))
	for i, item := range items {
		response[i] = *converter.ItemToResponse(&item)
	}

	return response, total, nil
}
