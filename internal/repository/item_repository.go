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
	base *BaseRepository[dbmodel.Items]
	Log  *logrus.Logger
}

func NewItemRepository(db *bun.DB, log *logrus.Logger) *ItemRepository {
	return &ItemRepository{
		base: NewBaseRepository[dbmodel.Items](db),
		Log:  log,
	}
}

func (r *ItemRepository) Create(ctx context.Context, tx bun.IDB, item *dbmodel.Items) (int64, error) {
	err := r.base.Insert(ctx, tx, item, "id")
	if err != nil {
		return 0, err
	}
	return item.ID, nil
}

func (r *ItemRepository) Update(ctx context.Context, tx bun.IDB, item *dbmodel.Items) error {
	return r.base.UpdateByPK(ctx, tx, item, "name", "sku", "currency", "stock", "updated_at")
}

func (r *ItemRepository) Get(ctx context.Context, tx bun.IDB, id int64) (*dbmodel.Items, error) {
	item := new(dbmodel.Items)
	err := r.base.FindOne(ctx, tx, item, func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("id = ?", id).Limit(1)
	})
	if err != nil {
		if apperror.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (r *ItemRepository) Delete(ctx context.Context, tx bun.IDB, id int64) error {
	return r.base.DeleteWhere(ctx, tx, "id = ?", id)
}

func (r *ItemRepository) Search(ctx context.Context, tx bun.IDB, search *dto.SearchItemRequest) ([]dbmodel.Items, int64, error) {
	var items []dbmodel.Items
	offset := (search.Page - 1) * search.Size

	applyFilter := func(q *bun.SelectQuery) *bun.SelectQuery {
		if name := search.Name; name != "" {
			pattern := "%" + name + "%"
			q = q.Where("name ILIKE ?", pattern)
		}

		if sku := search.SKU; sku != "" {
			pattern := "%" + sku + "%"
			q = q.Where("sku ILIKE ?", pattern)
		}

		if orderExpr := r.sortItem(search.Sort); orderExpr != "" {
			q = q.OrderExpr(orderExpr)
		}

		return q
	}

	count, err := r.base.Count(ctx, tx, applyFilter)
	if err != nil {
		return nil, 0, err
	}

	err = r.base.FindAll(ctx, tx, &items, func(q *bun.SelectQuery) *bun.SelectQuery {
		return applyFilter(q).Limit(search.Size).Offset(offset)
	})
	if err != nil {
		return nil, 0, err
	}

	return items, count, nil
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
