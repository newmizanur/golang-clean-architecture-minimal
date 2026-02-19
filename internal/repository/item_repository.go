package repository

import (
	"context"

	"golang-clean-architecture/ent"
	"golang-clean-architecture/ent/item"
	"golang-clean-architecture/internal/dto"

	"github.com/sirupsen/logrus"
)

type ItemRepository struct {
	Log *logrus.Logger
}

func NewItemRepository(log *logrus.Logger) *ItemRepository {
	return &ItemRepository{Log: log}
}

func (r *ItemRepository) Create(ctx context.Context, client *ent.Client, i *ent.Item) (*ent.Item, error) {
	created, err := client.Item.Create().
		SetName(i.Name).
		SetSku(i.Sku).
		SetCurrency(i.Currency).
		SetStock(i.Stock).
		SetCreatedAt(i.CreatedAt).
		SetUpdatedAt(i.UpdatedAt).
		Save(ctx)
	if err != nil {
		r.Log.WithError(err).Error("Failed to create item")
		return nil, err
	}
	r.Log.WithField("item_id", created.ID).Debug("Item created successfully")
	return created, nil
}

func (r *ItemRepository) FindById(ctx context.Context, client *ent.Client, id int) (*ent.Item, error) {
	i, err := client.Item.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return i, nil
}

func (r *ItemRepository) Update(ctx context.Context, client *ent.Client, i *ent.Item) error {
	_, err := client.Item.UpdateOneID(i.ID).
		SetName(i.Name).
		SetSku(i.Sku).
		SetCurrency(i.Currency).
		SetStock(i.Stock).
		SetUpdatedAt(i.UpdatedAt).
		Save(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("item_id", i.ID).Error("Failed to update item")
		return err
	}
	r.Log.WithField("item_id", i.ID).Debug("Item updated successfully")
	return nil
}

func (r *ItemRepository) Delete(ctx context.Context, client *ent.Client, id int) error {
	err := client.Item.DeleteOneID(id).Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("item_id", id).Error("Failed to delete item")
		return err
	}
	r.Log.WithField("item_id", id).Debug("Item deleted successfully")
	return nil
}

func (r *ItemRepository) Search(ctx context.Context, client *ent.Client, search *dto.SearchItemRequest) ([]*ent.Item, int64, error) {
	if search.Page < 1 {
		search.Page = 1
	}
	if search.Size < 1 {
		search.Size = 10
	}

	query := client.Item.Query()

	if search.Name != "" {
		query = query.Where(item.NameContainsFold(search.Name))
	}
	if search.SKU != "" {
		query = query.Where(item.SkuContainsFold(search.SKU))
	}

	if orderFunc := r.sortItem(search.Sort); orderFunc != nil {
		query = query.Order(orderFunc)
	}

	total, err := query.Count(ctx)
	if err != nil {
		r.Log.WithError(err).Error("Failed to count items")
		return nil, 0, err
	}

	offset := (search.Page - 1) * search.Size
	items, err := query.Limit(search.Size).Offset(offset).All(ctx)
	if err != nil {
		r.Log.WithError(err).Error("Failed to search items")
		return nil, 0, err
	}

	return items, int64(total), nil
}

var sortableItemCols = map[string]string{
	"name":      item.FieldName,
	"sku":       item.FieldSku,
	"stock":     item.FieldStock,
	"createdAt": item.FieldCreatedAt,
	"updatedAt": item.FieldUpdatedAt,
}

func (r *ItemRepository) sortItem(sort string) item.OrderOption {
	if sort == "" {
		return nil
	}

	desc := false
	if sort[0] == '-' {
		desc = true
		sort = sort[1:]
	}

	if sort == "" {
		return nil
	}

	column, ok := sortableItemCols[sort]
	if !ok {
		return nil
	}

	if desc {
		return item.OrderOption(ent.Desc(column))
	}
	return item.OrderOption(ent.Asc(column))
}
