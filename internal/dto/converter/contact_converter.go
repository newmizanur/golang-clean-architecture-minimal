package converter

import (
	"golang-clean-architecture/internal/dto"
	dbmodel "golang-clean-architecture/internal/persistence/model"
)

func ContactToResponse(contact *dbmodel.Contact) *dto.ContactResponse {
	return &dto.ContactResponse{
		ID:        contact.ID,
		FirstName: contact.FirstName,
		LastName:  stringValueOrEmpty(contact.LastName),
		Email:     stringValueOrEmpty(contact.Email),
		Phone:     stringValueOrEmpty(contact.Phone),
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}
}
