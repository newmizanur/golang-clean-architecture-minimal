package repository

import (
	"context"
	"golang-clean-architecture/internal/apperror"
	m "golang-clean-architecture/internal/persistence/model"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type AddressRepository struct {
	DB  *bun.DB
	Log *logrus.Logger
}

func NewAddressRepository(db *bun.DB, log *logrus.Logger) *AddressRepository {
	return &AddressRepository{
		DB:  db,
		Log: log,
	}
}

func (r *AddressRepository) dbConn(tx bun.IDB) bun.IDB {
	if tx != nil {
		return tx
	}
	return r.DB
}

func (r *AddressRepository) FindByIdAndContactId(ctx context.Context, tx bun.IDB, id string, contactId string) (*m.Address, error) {
	address := new(m.Address)
	err := r.dbConn(tx).NewSelect().
		Model(address).
		Where("id = ?", id).
		Where("contact_id = ?", contactId).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if apperror.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return address, nil
}

func (r *AddressRepository) FindAllByContactId(ctx context.Context, tx bun.IDB, contactId string) ([]m.Address, error) {
	var addresses []m.Address
	err := r.dbConn(tx).NewSelect().
		Model(&addresses).
		Where("contact_id = ?", contactId).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *AddressRepository) Create(ctx context.Context, tx bun.IDB, address *m.Address) error {
	_, err := r.dbConn(tx).NewInsert().Model(address).Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("address_id", address.ID).Error("Failed to create address")
		return err
	}
	r.Log.WithField("address_id", address.ID).Debug("Address created successfully")
	return nil
}

func (r *AddressRepository) Update(ctx context.Context, tx bun.IDB, address *m.Address) error {
	result, err := r.dbConn(tx).NewUpdate().
		Model(address).
		Column("street", "city", "province", "postal_code", "country", "updated_at").
		WherePK().
		Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("address_id", address.ID).Error("Failed to update address")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.Log.WithError(err).Error("Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		r.Log.WithField("address_id", address.ID).Warn("No address updated - address not found")
		return nil
	}

	r.Log.WithField("address_id", address.ID).Debug("Address updated successfully")
	return nil
}

func (r *AddressRepository) Delete(ctx context.Context, tx bun.IDB, address *m.Address) error {
	result, err := r.dbConn(tx).NewDelete().Model(address).WherePK().Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("address_id", address.ID).Error("Failed to delete address")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.Log.WithError(err).Error("Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		r.Log.WithField("address_id", address.ID).Warn("No address deleted - address not found")
		return nil
	}

	r.Log.WithField("address_id", address.ID).Debug("Address deleted successfully")
	return nil
}
