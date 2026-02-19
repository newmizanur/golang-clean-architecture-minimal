package repository

import (
	"context"

	"golang-clean-architecture/ent"
	"golang-clean-architecture/ent/user"

	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{Log: log}
}

func (r *UserRepository) CountById(ctx context.Context, client *ent.Client, id string) (int64, error) {
	count, err := client.User.Query().Where(user.ID(id)).Count(ctx)
	return int64(count), err
}

func (r *UserRepository) FindById(ctx context.Context, client *ent.Client, id string) (*ent.User, error) {
	u, err := client.User.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) Create(ctx context.Context, client *ent.Client, u *ent.User) error {
	_, err := client.User.Create().
		SetID(u.ID).
		SetName(u.Name).
		SetPassword(u.Password).
		SetCreatedAt(u.CreatedAt).
		SetUpdatedAt(u.UpdatedAt).
		Save(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("user_id", u.ID).Error("Failed to create user")
		return err
	}
	r.Log.WithField("user_id", u.ID).Debug("User created successfully")
	return nil
}

func (r *UserRepository) Update(ctx context.Context, client *ent.Client, u *ent.User) error {
	_, err := client.User.UpdateOneID(u.ID).
		SetName(u.Name).
		SetPassword(u.Password).
		SetUpdatedAt(u.UpdatedAt).
		Save(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("user_id", u.ID).Error("Failed to update user")
		return err
	}
	r.Log.WithField("user_id", u.ID).Debug("User updated successfully")
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, client *ent.Client, id string) error {
	err := client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("user_id", id).Error("Failed to delete user")
		return err
	}
	r.Log.WithField("user_id", id).Debug("User deleted successfully")
	return nil
}
