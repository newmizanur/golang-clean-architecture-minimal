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

func (r *AddressRepository) FindByIdAndContactId(ctx context.Context, tx bun.IDB, id string, contactId string) (*m.Addresses, error) {
	address := new(m.Addresses)
	err := r.dbConn(tx).NewSelect().
		Model(address).
		Where(m.AddressCols.ID+" = ?", id).
		Where(m.AddressCols.ContactID+" = ?", contactId).
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

func (r *AddressRepository) FindAllByContactId(ctx context.Context, tx bun.IDB, contactId string) ([]m.Addresses, error) {
	var addresses []m.Addresses
	err := r.dbConn(tx).NewSelect().
		Model(&addresses).
		Where(m.AddressCols.ContactID+" = ?", contactId).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *AddressRepository) Create(ctx context.Context, tx bun.IDB, address *m.Addresses) error {
	_, err := r.dbConn(tx).NewInsert().Model(address).Exec(ctx)
	return err
}

func (r *AddressRepository) Update(ctx context.Context, tx bun.IDB, address *m.Addresses) error {
	_, err := r.dbConn(tx).NewUpdate().
		Model(address).
		Column(m.AddressCols.Street, m.AddressCols.City, m.AddressCols.Province, m.AddressCols.PostalCode, m.AddressCols.Country, m.AddressCols.UpdatedAt).
		WherePK().
		Exec(ctx)
	return err
}

func (r *AddressRepository) Delete(ctx context.Context, tx bun.IDB, address *m.Addresses) error {
	_, err := r.dbConn(tx).NewDelete().Model(address).WherePK().Exec(ctx)
	return err
}
