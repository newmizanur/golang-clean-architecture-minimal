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
	_, err := r.dbConn(tx).NewInsert().Model(item).Returning(m.ItemCols.ID).Exec(ctx)
	if err != nil {
		return 0, err
	}
	return item.ID, nil
}

func (r *ItemRepository) Update(ctx context.Context, tx bun.IDB, item *m.Items) error {
	_, err := r.dbConn(tx).NewUpdate().
		Model(item).
		Column(m.ItemCols.Name, m.ItemCols.SKU, m.ItemCols.Currency, m.ItemCols.Stock, m.ItemCols.UpdatedAt).
		WherePK().
		Exec(ctx)
	return err
}

func (r *ItemRepository) Get(ctx context.Context, tx bun.IDB, id int64) (*m.Items, error) {
	item := new(m.Items)
	err := r.dbConn(tx).NewSelect().
		Model(item).
		Where(m.ItemCols.ID+" = ?", id).
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
		Model((*m.Items)(nil)).
		Where(m.ItemCols.ID+" = ?", id).
		Exec(ctx)
	return err
}

func (r *ItemRepository) Search(ctx context.Context, tx bun.IDB, search *dto.SearchItemRequest) ([]m.Items, int64, error) {
	var items []m.Items
	offset := (search.Page - 1) * search.Size

	query := r.dbConn(tx).NewSelect().Model(&items)

	if name := search.Name; name != "" {
		pattern := "%" + name + "%"
		query = query.Where(m.ItemCols.Name+" ILIKE ?", pattern)
	}

	if sku := search.SKU; sku != "" {
		pattern := "%" + sku + "%"
		query = query.Where(m.ItemCols.SKU+" ILIKE ?", pattern)
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
	"name":      m.ItemCols.Name,
	"sku":       m.ItemCols.SKU,
	"stock":     m.ItemCols.Stock,
	"createdAt": m.ItemCols.CreatedAt,
	"updatedAt": m.ItemCols.UpdatedAt,
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
