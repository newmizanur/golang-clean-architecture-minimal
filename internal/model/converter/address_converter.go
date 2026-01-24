package converter

import (
	dbmodel "golang-clean-architecture/internal/entity/db/model"
	"golang-clean-architecture/internal/model"
)

func AddressToResponse(address *dbmodel.Addresses) *model.AddressResponse {
	return &model.AddressResponse{
		ID:         address.ID,
		Street:     stringValueOrEmpty(address.Street),
		City:       stringValueOrEmpty(address.City),
		Province:   stringValueOrEmpty(address.Province),
		PostalCode: stringValueOrEmpty(address.PostalCode),
		Country:    stringValueOrEmpty(address.Country),
		CreatedAt:  address.CreatedAt,
		UpdatedAt:  address.UpdatedAt,
	}
}
