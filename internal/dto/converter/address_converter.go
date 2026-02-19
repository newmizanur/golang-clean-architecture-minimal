package converter

import (
	"golang-clean-architecture/ent"
	"golang-clean-architecture/internal/dto"
)

func AddressToResponse(address *ent.Address) *dto.AddressResponse {
	return &dto.AddressResponse{
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
