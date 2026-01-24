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

func (r *ContactRepository) FindByIdAndUserId(ctx context.Context, tx bun.IDB, id string, userId string) (*m.Contact, error) {
	contact := new(m.Contact)
	err := r.dbConn(tx).NewSelect().
		Model(contact).
		Where("id = ?", id).
		Where("user_id = ?", userId).
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

func (r *ContactRepository) Search(ctx context.Context, tx bun.IDB, request *dto.SearchContactRequest) ([]m.Contact, int64, error) {
	var contacts []m.Contact

	// Validate pagination parameters
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Size < 1 {
		request.Size = 10
	}

	offset := (request.Page - 1) * request.Size

	query := r.dbConn(tx).NewSelect().
		Model(&contacts).
		Where("user_id = ?", request.UserId)

	if name := request.Name; name != "" {
		pattern := "%" + name + "%"
		query = query.Where("(first_name ILIKE ? OR last_name ILIKE ?)", pattern, pattern)
	}

	if phone := request.Phone; phone != "" {
		pattern := "%" + phone + "%"
		query = query.Where("phone ILIKE ?", pattern)
	}

	if email := request.Email; email != "" {
		pattern := "%" + email + "%"
		query = query.Where("email ILIKE ?", pattern)
	}

	count, err := query.Count(ctx)
	if err != nil {
		r.Log.WithError(err).Error("Failed to count contacts")
		return nil, 0, err
	}

	err = query.Limit(request.Size).Offset(offset).Scan(ctx)
	if err != nil {
		r.Log.WithError(err).Error("Failed to search contacts")
		return nil, 0, err
	}

	return contacts, int64(count), nil
}

func (r *ContactRepository) Create(ctx context.Context, tx bun.IDB, contact *m.Contact) error {
	_, err := r.dbConn(tx).NewInsert().Model(contact).Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("contact_id", contact.ID).Error("Failed to create contact")
		return err
	}
	r.Log.WithField("contact_id", contact.ID).Debug("Contact created successfully")
	return nil
}

func (r *ContactRepository) Update(ctx context.Context, tx bun.IDB, contact *m.Contact) error {
	result, err := r.dbConn(tx).NewUpdate().
		Model(contact).
		OmitZero().
		WherePK().
		Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("contact_id", contact.ID).Error("Failed to update contact")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.Log.WithError(err).Error("Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		r.Log.WithField("contact_id", contact.ID).Warn("No contact updated - contact not found")
		return nil
	}

	r.Log.WithField("contact_id", contact.ID).Debug("Contact updated successfully")
	return nil
}

func (r *ContactRepository) Delete(ctx context.Context, tx bun.IDB, contact *m.Contact) error {
	result, err := r.dbConn(tx).NewDelete().Model(contact).WherePK().Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("contact_id", contact.ID).Error("Failed to delete contact")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.Log.WithError(err).Error("Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		r.Log.WithField("contact_id", contact.ID).Warn("No contact deleted - contact not found")
		return nil
	}

	r.Log.WithField("contact_id", contact.ID).Debug("Contact deleted successfully")
	return nil
}
