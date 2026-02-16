package repository

import (
	"context"

	"golang-clean-architecture/internal/apperror"
	"golang-clean-architecture/internal/dto"
	m "golang-clean-architecture/internal/persistence/model"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type ItemRepository struct {
	DB  *bun.DB
	Log *logrus.Logger
}

func NewItemRepository(db *bun.DB, log *logrus.Logger) *ItemRepository {
	return &ItemRepository{
		DB:  db,
		Log: log,
	}
}

func (r *ItemRepository) dbConn(tx bun.IDB) bun.IDB {
	if tx != nil {
		return tx
	}
	return r.DB
}

func (r *ItemRepository) Create(ctx context.Context, tx bun.IDB, item *m.Items) (int64, error) {
	_, err := r.dbConn(tx).NewInsert().Model(item).Returning("id").Exec(ctx)
	if err != nil {
		r.Log.WithError(err).Error("Failed to create item")
		return 0, err
	}
	r.Log.WithField("item_id", item.ID).Info("Item created successfully")
	return item.ID, nil
}

func (r *ItemRepository) Update(ctx context.Context, tx bun.IDB, item *m.Items) error {
	result, err := r.dbConn(tx).NewUpdate().
		Model(item).
		OmitZero().
		WherePK().
		Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("item_id", item.ID).Error("Failed to update item")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.Log.WithError(err).Error("Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		r.Log.WithField("item_id", item.ID).Warn("No item updated - item not found")
		return nil
	}

	r.Log.WithField("item_id", item.ID).Info("Item updated successfully")
	return nil
}

func (r *ItemRepository) FindById(ctx context.Context, tx bun.IDB, id int64) (*m.Items, error) {
	item := new(m.Items)
	err := r.dbConn(tx).NewSelect().
		Model(item).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if apperror.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (r *ItemRepository) Delete(ctx context.Context, tx bun.IDB, item *m.Items) error {
	result, err := r.dbConn(tx).NewDelete().Model(item).WherePK().Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("item_id", item.ID).Error("Failed to delete item")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.Log.WithError(err).Error("Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		r.Log.WithField("item_id", item.ID).Warn("No item deleted - item not found")
		return nil
	}

	r.Log.WithField("item_id", item.ID).Info("Item deleted successfully")
	return nil
}

func (r *ItemRepository) Search(ctx context.Context, tx bun.IDB, search *dto.SearchItemRequest) ([]m.Items, int64, error) {
	var items []m.Items

	// Validate pagination parameters
	if search.Page < 1 {
		search.Page = 1
	}
	if search.Size < 1 {
		search.Size = 10
	}

	offset := (search.Page - 1) * search.Size

	query := r.dbConn(tx).NewSelect().Model(&items)

	if name := search.Name; name != "" {
		pattern := "%" + name + "%"
		query = query.Where("name ILIKE ?", pattern)
	}

	if sku := search.SKU; sku != "" {
		pattern := "%" + sku + "%"
		query = query.Where("sku ILIKE ?", pattern)
	}

	if orderExpr := r.sortItem(search.Sort); orderExpr != "" {
		query = query.OrderExpr(orderExpr)
	}

	count, err := query.Count(ctx)
	if err != nil {
		r.Log.WithError(err).Error("Failed to count items")
		return nil, 0, err
	}

	err = query.Limit(search.Size).Offset(offset).Scan(ctx)
	if err != nil {
		r.Log.WithError(err).Error("Failed to search items")
		return nil, 0, err
	}

	return items, int64(count), nil
}

var sortableItemCols = map[string]string{
	"name":      "name",
	"sku":       "sku",
	"stock":     "stock",
	"createdAt": "created_at",
	"updatedAt": "updated_at",
}

func (r *ItemRepository) sortItem(sort string) string {
	if sort == "" {
		return ""
	}

	order := "ASC"
	if sort[0] == '-' {
		order = "DESC"
		sort = sort[1:]
	}

	if sort == "" {
		return ""
	}

	column, ok := sortableItemCols[sort]
	if !ok {
		return ""
	}

	return column + " " + order
}
