package repository

import (
	"context"
	"golang-clean-architecture/internal/apperror"
	dbmodel "golang-clean-architecture/internal/persistence/model"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	base *BaseRepository[dbmodel.Users]
	Log  *logrus.Logger
}

func NewUserRepository(db *bun.DB, log *logrus.Logger) *UserRepository {
	return &UserRepository{
		base: NewBaseRepository[dbmodel.Users](db),
		Log:  log,
	}
}

func (r *UserRepository) CountById(ctx context.Context, tx bun.IDB, id string) (int64, error) {
	return r.base.Count(ctx, tx, func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("id = ?", id)
	})
}

func (r *UserRepository) FindById(ctx context.Context, tx bun.IDB, id string) (*dbmodel.Users, error) {
	user := new(dbmodel.Users)
	err := r.base.FindOne(ctx, tx, user, func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("id = ?", id).Limit(1)
	})
	if err != nil {
		if apperror.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, tx bun.IDB, user *dbmodel.Users) error {
	return r.base.Insert(ctx, tx, user)
}

func (r *UserRepository) Update(ctx context.Context, tx bun.IDB, user *dbmodel.Users) error {
	return r.base.UpdateByPK(ctx, tx, user, "password", "name", "updated_at")
}

func (r *UserRepository) Delete(ctx context.Context, tx bun.IDB, user *dbmodel.Users) error {
	return r.base.DeleteByPK(ctx, tx, user)
}
