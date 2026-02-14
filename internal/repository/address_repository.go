package repository

import (
	"context"
	"golang-clean-architecture/internal/apperror"
	dbmodel "golang-clean-architecture/internal/persistence/model"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type AddressRepository struct {
	base *BaseRepository[dbmodel.Addresses]
	Log  *logrus.Logger
}

func NewAddressRepository(db *bun.DB, log *logrus.Logger) *AddressRepository {
	return &AddressRepository{
		base: NewBaseRepository[dbmodel.Addresses](db),
		Log:  log,
	}
}

func (r *AddressRepository) FindByIdAndContactId(ctx context.Context, tx bun.IDB, id string, contactId string) (*dbmodel.Addresses, error) {
	address := new(dbmodel.Addresses)
	err := r.base.FindOne(ctx, tx, address, func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("id = ?", id).Where("contact_id = ?", contactId).Limit(1)
	})
	if err != nil {
		if apperror.IsNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return address, nil
}

func (r *AddressRepository) FindAllByContactId(ctx context.Context, tx bun.IDB, contactId string) ([]dbmodel.Addresses, error) {
	var addresses []dbmodel.Addresses
	err := r.base.FindAll(ctx, tx, &addresses, func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("contact_id = ?", contactId)
	})
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *AddressRepository) Create(ctx context.Context, tx bun.IDB, address *dbmodel.Addresses) error {
	return r.base.Insert(ctx, tx, address)
}

func (r *AddressRepository) Update(ctx context.Context, tx bun.IDB, address *dbmodel.Addresses) error {
	return r.base.UpdateByPK(ctx, tx, address, "street", "city", "province", "postal_code", "country", "updated_at")
}

func (r *AddressRepository) Delete(ctx context.Context, tx bun.IDB, address *dbmodel.Addresses) error {
	return r.base.DeleteByPK(ctx, tx, address)
}
