package repository

import (
	"context"

	"golang-clean-architecture/internal/apperror"
	"golang-clean-architecture/internal/dto"
	dbmodel "golang-clean-architecture/internal/persistence/model"

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

func (r *ItemRepository) Create(ctx context.Context, tx bun.IDB, item *dbmodel.Items) (int64, error) {
	_, err := r.dbConn(tx).NewInsert().Model(item).Returning("id").Exec(ctx)
	if err != nil {
		return 0, err
	}
	return item.ID, nil
}

func (r *ItemRepository) Update(ctx context.Context, tx bun.IDB, item *dbmodel.Items) error {
	_, err := r.dbConn(tx).NewUpdate().
		Model(item).
		Column("name", "sku", "currency", "stock", "updated_at").
		WherePK().
		Exec(ctx)
	return err
}

func (r *ItemRepository) Get(ctx context.Context, tx bun.IDB, id int64) (*dbmodel.Items, error) {
	item := new(dbmodel.Items)
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

func (r *ItemRepository) Delete(ctx context.Context, tx bun.IDB, id int64) error {
	_, err := r.dbConn(tx).NewDelete().
		Model((*dbmodel.Items)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *ItemRepository) Search(ctx context.Context, tx bun.IDB, search *dto.SearchItemRequest) ([]dbmodel.Items, int64, error) {
	var items []dbmodel.Items
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
		return nil, 0, err
	}

	err = query.Limit(search.Size).Offset(offset).Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return items, int64(count), nil
}

func (r *ItemRepository) sortItem(sort string) string {
	switch sort {
	case "name":
		return "name ASC"
	case "-name":
		return "name DESC"
	case "sku":
		return "sku ASC"
	case "-sku":
		return "sku DESC"
	case "stock":
		return "stock ASC"
	case "-stock":
		return "stock DESC"
	case "createdAt":
		return "created_at ASC"
	case "-createdAt":
		return "created_at DESC"
	case "updatedAt":
		return "updated_at ASC"
	case "-updatedAt":
		return "updated_at DESC"
	default:
		return ""
	}
}
