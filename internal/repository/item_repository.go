package repository

import (
	"context"
	"database/sql"

	"golang-clean-architecture/internal/apperror"
	dbmodel "golang-clean-architecture/internal/entity/db/model"
	t "golang-clean-architecture/internal/entity/db/table"
	"golang-clean-architecture/internal/model"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/sirupsen/logrus"
)

type ItemRepository struct {
	DB  *sql.DB
	Log *logrus.Logger
}

func NewItemRepository(db *sql.DB, log *logrus.Logger) *ItemRepository {
	return &ItemRepository{
		DB:  db,
		Log: log,
	}
}

func (r *ItemRepository) Create(ctx context.Context, tx *sql.Tx, item *dbmodel.Items) (int64, error) {
	stmt := t.Items.INSERT(
		t.Items.Name,
		t.Items.Sku,
		t.Items.Currency,
		t.Items.Stock,
		t.Items.CreatedAt,
		t.Items.UpdatedAt,
	).MODEL(item)

	var db qrm.DB = r.DB
	if tx != nil {
		db = tx
	}

	result, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *ItemRepository) Update(ctx context.Context, tx *sql.Tx, item *dbmodel.Items) error {
	stmt := t.Items.UPDATE(
		t.Items.Name,
		t.Items.Sku,
		t.Items.Currency,
		t.Items.Stock,
	).MODEL(item).WHERE(t.Items.ID.EQ(mysql.Uint64(item.ID)))

	var db qrm.DB = r.DB
	if tx != nil {
		db = tx
	}

	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (r *ItemRepository) Get(ctx context.Context, tx *sql.Tx, id int64) (*dbmodel.Items, error) {
	stmt := t.Items.SELECT(
		t.Items.AllColumns,
	).WHERE(t.Items.ID.EQ(mysql.Int64(id))).LIMIT(1)

	var db qrm.DB = r.DB
	if tx != nil {
		db = tx
		stmt = stmt.FOR(mysql.UPDATE())
	}

	item := new(dbmodel.Items)
	if err := stmt.QueryContext(ctx, db, item); err != nil {
		if apperror.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (r *ItemRepository) Delete(ctx context.Context, tx *sql.Tx, id int64) error {
	stmt := t.Items.DELETE().WHERE(t.Items.ID.EQ(mysql.Int64(id))).LIMIT(1)

	var db qrm.DB = r.DB
	if tx != nil {
		db = tx
	}

	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (r *ItemRepository) Search(ctx context.Context, tx *sql.Tx, search *model.SearchItemRequest) ([]dbmodel.Items, int64, error) {
	var items []dbmodel.Items
	filter := r.filterItem(search)
	offset := (search.Page - 1) * search.Size

	stmt := t.Items.SELECT(
		t.Items.AllColumns,
	).WHERE(filter).LIMIT(int64(search.Size)).OFFSET(int64(offset))
	if orderBy := r.sortItem(search.Sort); len(orderBy) > 0 {
		stmt = stmt.ORDER_BY(orderBy...)
	}

	var db qrm.DB = r.DB
	if tx != nil {
		db = tx
	}

	//For data
	if err := stmt.QueryContext(ctx, db, &items); err != nil {
		if apperror.IsNoRows(err) {
			return items, 0, nil
		}
		return nil, 0, err
	}

	var result struct {
		Total int64
	}

	countStmt := t.Items.SELECT(
		mysql.COUNT(t.Items.ID).AS("total")).WHERE(filter)
	if err := countStmt.QueryContext(ctx, db, &result); err != nil {
		if apperror.IsNoRows(err) {
			return items, 0, nil
		}
		return nil, 0, err
	}

	return items, result.Total, nil
}

func (r *ItemRepository) filterItem(search *model.SearchItemRequest) mysql.BoolExpression {
	condition := mysql.Bool(true)

	if name := search.Name; name != "" {
		pattern := "%" + name + "%"
		condition = condition.AND(t.Items.Name.LIKE(mysql.String(pattern)))
	}

	if sku := search.SKU; sku != "" {
		pattern := "%" + sku + "%"
		condition = condition.AND(t.Items.Sku.LIKE(mysql.String(pattern)))
	}

	return condition
}

func (r *ItemRepository) sortItem(sort string) []mysql.OrderByClause {
	switch sort {
	case "name":
		return []mysql.OrderByClause{t.Items.Name.ASC()}
	case "-name":
		return []mysql.OrderByClause{t.Items.Name.DESC()}
	case "sku":
		return []mysql.OrderByClause{t.Items.Sku.ASC()}
	case "-sku":
		return []mysql.OrderByClause{t.Items.Sku.DESC()}
	case "stock":
		return []mysql.OrderByClause{t.Items.Stock.ASC()}
	case "-stock":
		return []mysql.OrderByClause{t.Items.Stock.DESC()}
	case "createdAt":
		return []mysql.OrderByClause{t.Items.CreatedAt.ASC()}
	case "-createdAt":
		return []mysql.OrderByClause{t.Items.CreatedAt.DESC()}
	case "updatedAt":
		return []mysql.OrderByClause{t.Items.UpdatedAt.ASC()}
	case "-updatedAt":
		return []mysql.OrderByClause{t.Items.UpdatedAt.DESC()}
	default:
		return nil
	}
}
