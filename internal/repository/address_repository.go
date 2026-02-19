package repository

import (
	"context"

	"golang-clean-architecture/ent"
	"golang-clean-architecture/ent/address"

	"github.com/sirupsen/logrus"
)

type AddressRepository struct {
	Log *logrus.Logger
}

func NewAddressRepository(log *logrus.Logger) *AddressRepository {
	return &AddressRepository{Log: log}
}

func (r *AddressRepository) FindByIdAndContactId(ctx context.Context, client *ent.Client, id string, contactId string) (*ent.Address, error) {
	a, err := client.Address.Query().
		Where(address.ID(id), address.ContactID(contactId)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

func (r *AddressRepository) FindAllByContactId(ctx context.Context, client *ent.Client, contactId string) ([]*ent.Address, error) {
	addresses, err := client.Address.Query().
		Where(address.ContactID(contactId)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *AddressRepository) Create(ctx context.Context, client *ent.Client, a *ent.Address) error {
	builder := client.Address.Create().
		SetID(a.ID).
		SetContactID(a.ContactID).
		SetCreatedAt(a.CreatedAt).
		SetUpdatedAt(a.UpdatedAt)

	if a.Street != nil {
		builder = builder.SetNillableStreet(a.Street)
	}
	if a.City != nil {
		builder = builder.SetNillableCity(a.City)
	}
	if a.Province != nil {
		builder = builder.SetNillableProvince(a.Province)
	}
	if a.PostalCode != nil {
		builder = builder.SetNillablePostalCode(a.PostalCode)
	}
	if a.Country != nil {
		builder = builder.SetNillableCountry(a.Country)
	}

	_, err := builder.Save(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("address_id", a.ID).Error("Failed to create address")
		return err
	}
	r.Log.WithField("address_id", a.ID).Debug("Address created successfully")
	return nil
}

func (r *AddressRepository) Update(ctx context.Context, client *ent.Client, a *ent.Address) error {
	builder := client.Address.UpdateOneID(a.ID).
		SetUpdatedAt(a.UpdatedAt)

	if a.Street != nil {
		builder = builder.SetStreet(*a.Street)
	} else {
		builder = builder.ClearStreet()
	}
	if a.City != nil {
		builder = builder.SetCity(*a.City)
	} else {
		builder = builder.ClearCity()
	}
	if a.Province != nil {
		builder = builder.SetProvince(*a.Province)
	} else {
		builder = builder.ClearProvince()
	}
	if a.PostalCode != nil {
		builder = builder.SetPostalCode(*a.PostalCode)
	} else {
		builder = builder.ClearPostalCode()
	}
	if a.Country != nil {
		builder = builder.SetCountry(*a.Country)
	} else {
		builder = builder.ClearCountry()
	}

	_, err := builder.Save(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("address_id", a.ID).Error("Failed to update address")
		return err
	}
	r.Log.WithField("address_id", a.ID).Debug("Address updated successfully")
	return nil
}

func (r *AddressRepository) Delete(ctx context.Context, client *ent.Client, id string) error {
	err := client.Address.DeleteOneID(id).Exec(ctx)
	if err != nil {
		r.Log.WithError(err).WithField("address_id", id).Error("Failed to delete address")
		return err
	}
	r.Log.WithField("address_id", id).Debug("Address deleted successfully")
	return nil
}
