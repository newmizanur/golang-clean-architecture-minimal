package repository

import (
	"context"

	"github.com/uptrace/bun"
)

type BaseRepository[T any] struct {
	DB *bun.DB
}

func NewBaseRepository[T any](db *bun.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

func (r *BaseRepository[T]) dbConn(tx bun.IDB) bun.IDB {
	if tx != nil {
		return tx
	}
	return r.DB
}

func (r *BaseRepository[T]) Insert(ctx context.Context, tx bun.IDB, model *T, returning ...string) error {
	query := r.dbConn(tx).NewInsert().Model(model)
	for _, column := range returning {
		query = query.Returning(column)
	}
	_, err := query.Exec(ctx)
	return err
}

func (r *BaseRepository[T]) UpdateByPK(ctx context.Context, tx bun.IDB, model *T, columns ...string) error {
	query := r.dbConn(tx).NewUpdate().Model(model).WherePK()
	if len(columns) > 0 {
		query = query.Column(columns...)
	}
	_, err := query.Exec(ctx)
	return err
}

func (r *BaseRepository[T]) DeleteByPK(ctx context.Context, tx bun.IDB, model *T) error {
	_, err := r.dbConn(tx).NewDelete().Model(model).WherePK().Exec(ctx)
	return err
}

func (r *BaseRepository[T]) DeleteWhere(ctx context.Context, tx bun.IDB, where string, args ...any) error {
	_, err := r.dbConn(tx).NewDelete().Model((*T)(nil)).Where(where, args...).Exec(ctx)
	return err
}

func (r *BaseRepository[T]) FindOne(ctx context.Context, tx bun.IDB, model *T, queryFn func(*bun.SelectQuery) *bun.SelectQuery) error {
	query := r.dbConn(tx).NewSelect().Model(model)
	if queryFn != nil {
		query = queryFn(query)
	}
	return query.Scan(ctx)
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, tx bun.IDB, models *[]T, queryFn func(*bun.SelectQuery) *bun.SelectQuery) error {
	query := r.dbConn(tx).NewSelect().Model(models)
	if queryFn != nil {
		query = queryFn(query)
	}
	return query.Scan(ctx)
}

func (r *BaseRepository[T]) Count(ctx context.Context, tx bun.IDB, queryFn func(*bun.SelectQuery) *bun.SelectQuery) (int64, error) {
	query := r.dbConn(tx).NewSelect().Model((*T)(nil))
	if queryFn != nil {
		query = queryFn(query)
	}
	count, err := query.Count(ctx)
	return int64(count), err
}
