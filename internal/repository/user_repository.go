package repository

import (
	"context"
	"golang-clean-architecture/internal/apperror"
	m "golang-clean-architecture/internal/persistence/model"

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
		Model((*m.Users)(nil)).
		Where("id = ?", id).
		Count(ctx)
	return int64(count), err
}

func (r *UserRepository) FindById(ctx context.Context, tx bun.IDB, id string) (*m.Users, error) {
	user := new(m.Users)
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

func (r *UserRepository) Create(ctx context.Context, tx bun.IDB, user *m.Users) error {
	_, err := r.dbConn(tx).NewInsert().Model(user).Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("user_id", user.ID).Error("Failed to create user")
		return err
	}
	r.Log.WithField("user_id", user.ID).Info("User created successfully")
	return nil
}

func (r *UserRepository) Update(ctx context.Context, tx bun.IDB, user *m.Users) error {
	result, err := r.dbConn(tx).NewUpdate().
		Model(user).
		OmitZero().
		WherePK().
		Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("user_id", user.ID).Error("Failed to update user")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.Log.WithError(err).Error("Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		r.Log.WithField("user_id", user.ID).Warn("No user updated - user not found")
		return nil
	}

	r.Log.WithField("user_id", user.ID).Info("User updated successfully")
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, tx bun.IDB, user *m.Users) error {
	result, err := r.dbConn(tx).NewDelete().Model(user).WherePK().Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("user_id", user.ID).Error("Failed to delete user")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.Log.WithError(err).Error("Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		r.Log.WithField("user_id", user.ID).Warn("No user deleted - user not found")
		return nil
	}

	r.Log.WithField("user_id", user.ID).Info("User deleted successfully")
	return nil
}
