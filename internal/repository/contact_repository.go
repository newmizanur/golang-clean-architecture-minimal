package repository

import (
	"context"
	"golang-clean-architecture/internal/apperror"
	"golang-clean-architecture/internal/dto"
	dbmodel "golang-clean-architecture/internal/persistence/model"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type ContactRepository struct {
	base *BaseRepository[dbmodel.Contacts]
	Log  *logrus.Logger
}

func NewContactRepository(db *bun.DB, log *logrus.Logger) *ContactRepository {
	return &ContactRepository{
		base: NewBaseRepository[dbmodel.Contacts](db),
		Log:  log,
	}
}

func (r *ContactRepository) FindByIdAndUserId(ctx context.Context, tx bun.IDB, id string, userId string) (*dbmodel.Contacts, error) {
	contact := new(dbmodel.Contacts)
	err := r.base.FindOne(ctx, tx, contact, func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("id = ?", id).Where("user_id = ?", userId).Limit(1)
	})
	if err != nil {
		if apperror.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return contact, nil
}

func (r *ContactRepository) Search(ctx context.Context, tx bun.IDB, request *dto.SearchContactRequest) ([]dbmodel.Contacts, int64, error) {
	var contacts []dbmodel.Contacts
	offset := (request.Page - 1) * request.Size

	applyFilter := func(q *bun.SelectQuery) *bun.SelectQuery {
		q = q.Where("user_id = ?", request.UserId)

		if name := request.Name; name != "" {
			pattern := "%" + name + "%"
			q = q.Where("(first_name ILIKE ? OR last_name ILIKE ?)", pattern, pattern)
		}

		if phone := request.Phone; phone != "" {
			pattern := "%" + phone + "%"
			q = q.Where("phone ILIKE ?", pattern)
		}

		if email := request.Email; email != "" {
			pattern := "%" + email + "%"
			q = q.Where("email ILIKE ?", pattern)
		}

		return q
	}

	count, err := r.base.Count(ctx, tx, applyFilter)
	if err != nil {
		return nil, 0, err
	}

	err = r.base.FindAll(ctx, tx, &contacts, func(q *bun.SelectQuery) *bun.SelectQuery {
		return applyFilter(q).Limit(request.Size).Offset(offset)
	})
	if err != nil {
		return nil, 0, err
	}

	return contacts, count, nil
}

func (r *ContactRepository) Create(ctx context.Context, tx bun.IDB, contact *dbmodel.Contacts) error {
	return r.base.Insert(ctx, tx, contact)
}

func (r *ContactRepository) Update(ctx context.Context, tx bun.IDB, contact *dbmodel.Contacts) error {
	return r.base.UpdateByPK(ctx, tx, contact, "first_name", "last_name", "email", "phone", "updated_at")
}

func (r *ContactRepository) Delete(ctx context.Context, tx bun.IDB, contact *dbmodel.Contacts) error {
	return r.base.DeleteByPK(ctx, tx, contact)
}
