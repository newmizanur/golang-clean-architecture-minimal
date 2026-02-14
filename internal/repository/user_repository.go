package repository

import (
	"context"
	"golang-clean-architecture/internal/apperror"
	dbmodel "golang-clean-architecture/internal/persistence/model"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	DB  *bun.DB
	Log *logrus.Logger
}

func NewUserRepository(db *bun.DB, log *logrus.Logger) *UserRepository {
	return &UserRepository{
		DB:  db,
		Log: log,
	}
}

func (r *UserRepository) dbConn(tx bun.IDB) bun.IDB {
	if tx != nil {
		return tx
	}
	return r.DB
}

func (r *UserRepository) CountById(ctx context.Context, tx bun.IDB, id string) (int64, error) {
	count, err := r.dbConn(tx).NewSelect().
		Model((*dbmodel.Users)(nil)).
		Where("id = ?", id).
		Count(ctx)
	return int64(count), err
}

func (r *UserRepository) FindById(ctx context.Context, tx bun.IDB, id string) (*dbmodel.Users, error) {
	user := new(dbmodel.Users)
	err := r.dbConn(tx).NewSelect().
		Model(user).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if apperror.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, tx bun.IDB, user *dbmodel.Users) error {
	_, err := r.dbConn(tx).NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *UserRepository) Update(ctx context.Context, tx bun.IDB, user *dbmodel.Users) error {
	_, err := r.dbConn(tx).NewUpdate().
		Model(user).
		Column("password", "name", "updated_at").
		WherePK().
		Exec(ctx)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, tx bun.IDB, user *dbmodel.Users) error {
	_, err := r.dbConn(tx).NewDelete().Model(user).WherePK().Exec(ctx)
	return err
}
