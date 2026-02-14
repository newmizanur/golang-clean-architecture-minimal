package repository

import (
	"context"
	"golang-clean-architecture/internal/apperror"
	"golang-clean-architecture/internal/dto"
	m "golang-clean-architecture/internal/persistence/model"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type ContactRepository struct {
	DB  *bun.DB
	Log *logrus.Logger
}

func NewContactRepository(db *bun.DB, log *logrus.Logger) *ContactRepository {
	return &ContactRepository{
		DB:  db,
		Log: log,
	}
}

func (r *ContactRepository) dbConn(tx bun.IDB) bun.IDB {
	if tx != nil {
		return tx
	}
	return r.DB
}

func (r *ContactRepository) FindByIdAndUserId(ctx context.Context, tx bun.IDB, id string, userId string) (*m.Contacts, error) {
	contact := new(m.Contacts)
	err := r.dbConn(tx).NewSelect().
		Model(contact).
		Where(m.ContactCols.ID+" = ?", id).
		Where(m.ContactCols.UserID+" = ?", userId).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if apperror.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return contact, nil
}

func (r *ContactRepository) Search(ctx context.Context, tx bun.IDB, request *dto.SearchContactRequest) ([]m.Contacts, int64, error) {
	var contacts []m.Contacts
	offset := (request.Page - 1) * request.Size

	query := r.dbConn(tx).NewSelect().
		Model(&contacts).
		Where(m.ContactCols.UserID+" = ?", request.UserId)

	if name := request.Name; name != "" {
		pattern := "%" + name + "%"
		query = query.Where("("+m.ContactCols.FirstName+" ILIKE ? OR "+m.ContactCols.LastName+" ILIKE ?)", pattern, pattern)
	}

	if phone := request.Phone; phone != "" {
		pattern := "%" + phone + "%"
		query = query.Where(m.ContactCols.Phone+" ILIKE ?", pattern)
	}

	if email := request.Email; email != "" {
		pattern := "%" + email + "%"
		query = query.Where(m.ContactCols.Email+" ILIKE ?", pattern)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(request.Size).Offset(offset).Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return contacts, int64(count), nil
}

func (r *ContactRepository) Create(ctx context.Context, tx bun.IDB, contact *m.Contacts) error {
	_, err := r.dbConn(tx).NewInsert().Model(contact).Exec(ctx)
	return err
}

func (r *ContactRepository) Update(ctx context.Context, tx bun.IDB, contact *m.Contacts) error {
	_, err := r.dbConn(tx).NewUpdate().
		Model(contact).
		Column(m.ContactCols.FirstName, m.ContactCols.LastName, m.ContactCols.Email, m.ContactCols.Phone, m.ContactCols.UpdatedAt).
		WherePK().
		Exec(ctx)
	return err
}

func (r *ContactRepository) Delete(ctx context.Context, tx bun.IDB, contact *m.Contacts) error {
	_, err := r.dbConn(tx).NewDelete().Model(contact).WherePK().Exec(ctx)
	return err
}
