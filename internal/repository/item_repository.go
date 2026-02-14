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
	_, err := r.dbConn(tx).NewInsert().Model(item).Returning(dbmodel.ItemCols.ID).Exec(ctx)
	if err != nil {
		return 0, err
	}
	return item.ID, nil
}

func (r *ItemRepository) Update(ctx context.Context, tx bun.IDB, item *dbmodel.Items) error {
	_, err := r.dbConn(tx).NewUpdate().
		Model(item).
		Column(dbmodel.ItemCols.Name, dbmodel.ItemCols.SKU, dbmodel.ItemCols.Currency, dbmodel.ItemCols.Stock, dbmodel.ItemCols.UpdatedAt).
		WherePK().
		Exec(ctx)
	return err
}

func (r *ItemRepository) Get(ctx context.Context, tx bun.IDB, id int64) (*dbmodel.Items, error) {
	item := new(dbmodel.Items)
	err := r.dbConn(tx).NewSelect().
		Model(item).
		Where(dbmodel.ItemCols.ID+" = ?", id).
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
		Where(dbmodel.ItemCols.ID+" = ?", id).
		Exec(ctx)
	return err
}

func (r *ItemRepository) Search(ctx context.Context, tx bun.IDB, search *dto.SearchItemRequest) ([]dbmodel.Items, int64, error) {
	var items []dbmodel.Items
	offset := (search.Page - 1) * search.Size

	query := r.dbConn(tx).NewSelect().Model(&items)

	if name := search.Name; name != "" {
		pattern := "%" + name + "%"
		query = query.Where(dbmodel.ItemCols.Name+" ILIKE ?", pattern)
	}

	if sku := search.SKU; sku != "" {
		pattern := "%" + sku + "%"
		query = query.Where(dbmodel.ItemCols.SKU+" ILIKE ?", pattern)
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

var sortableItemCols = map[string]string{
	"name":      dbmodel.ItemCols.Name,
	"sku":       dbmodel.ItemCols.SKU,
	"stock":     dbmodel.ItemCols.Stock,
	"createdAt": dbmodel.ItemCols.CreatedAt,
	"updatedAt": dbmodel.ItemCols.UpdatedAt,
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
