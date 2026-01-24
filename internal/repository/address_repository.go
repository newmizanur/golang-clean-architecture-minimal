package repository

import (
	"context"
	"database/sql"
	"golang-clean-architecture/internal/apperror"
	dbmodel "golang-clean-architecture/internal/entity/db/model"
	t "golang-clean-architecture/internal/entity/db/table"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/sirupsen/logrus"
)

type AddressRepository struct {
	DB  *sql.DB
	Log *logrus.Logger
}

func NewAddressRepository(db *sql.DB, log *logrus.Logger) *AddressRepository {
	return &AddressRepository{
		DB:  db,
		Log: log,
	}
}

func (r *AddressRepository) FindByIdAndContactId(ctx context.Context, tx *sql.Tx, id string, contactId string) (*dbmodel.Addresses, error) {
	stmt := mysql.SELECT(t.Addresses.AllColumns).
		FROM(t.Addresses).
		WHERE(
			t.Addresses.ID.EQ(mysql.String(id)).
				AND(t.Addresses.ContactID.EQ(mysql.String(contactId))),
		).
		LIMIT(1)
	db := qrm.Queryable(r.DB)
	if tx != nil {
		db = tx
		stmt = stmt.FOR(mysql.UPDATE())
	}
	address := new(dbmodel.Addresses)
	if err := stmt.QueryContext(ctx, db, address); err != nil {
		if apperror.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return address, nil
}

func (r *AddressRepository) FindAllByContactId(ctx context.Context, tx *sql.Tx, contactId string) ([]dbmodel.Addresses, error) {
	var addresses []dbmodel.Addresses
	stmt := mysql.SELECT(t.Addresses.AllColumns).
		FROM(t.Addresses).
		WHERE(t.Addresses.ContactID.EQ(mysql.String(contactId)))
	db := qrm.Queryable(r.DB)
	if tx != nil {
		db = tx
		stmt = stmt.FOR(mysql.UPDATE())
	}
	if err := stmt.QueryContext(ctx, db, &addresses); err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *AddressRepository) Create(ctx context.Context, tx *sql.Tx, address *dbmodel.Addresses) error {
	stmt := t.Addresses.INSERT(t.Addresses.ID,
		t.Addresses.ContactID,
		t.Addresses.Street,
		t.Addresses.City,
		t.Addresses.Province,
		t.Addresses.PostalCode,
		t.Addresses.Country,
		t.Addresses.CreatedAt,
		t.Addresses.UpdatedAt).MODEL(address)
	db := qrm.Executable(r.DB)
	if tx != nil {
		db = tx
	}
	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (r *AddressRepository) Update(ctx context.Context, tx *sql.Tx, address *dbmodel.Addresses) error {
	stmt := t.Addresses.UPDATE(
		t.Addresses.Street,
		t.Addresses.City,
		t.Addresses.Province,
		t.Addresses.PostalCode,
		t.Addresses.Country,
		t.Addresses.CreatedAt,
		t.Addresses.UpdatedAt,
	).
		SET(
			t.Addresses.Street.SET(stringExprOrNull(address.Street)),
			t.Addresses.City.SET(stringExprOrNull(address.City)),
			t.Addresses.Province.SET(stringExprOrNull(address.Province)),
			t.Addresses.PostalCode.SET(stringExprOrNull(address.PostalCode)),
			t.Addresses.Country.SET(stringExprOrNull(address.Country)),
			t.Addresses.UpdatedAt.SET(mysql.Int(address.UpdatedAt)),
		).
		WHERE(t.Addresses.ID.EQ(mysql.String(address.ID)))
	db := qrm.Executable(r.DB)
	if tx != nil {
		db = tx
	}
	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (r *AddressRepository) Delete(ctx context.Context, tx *sql.Tx, address *dbmodel.Addresses) error {
	stmt := t.Addresses.DELETE().
		WHERE(t.Addresses.ID.EQ(mysql.String(address.ID)))
	db := qrm.Executable(r.DB)
	if tx != nil {
		db = tx
	}
	_, err := stmt.ExecContext(ctx, db)
	return err
}
