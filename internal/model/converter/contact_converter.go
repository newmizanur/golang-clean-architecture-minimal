package converter

import (
	dbmodel "golang-clean-architecture/internal/entity/db/model"
	"golang-clean-architecture/internal/model"
)

func ContactToResponse(contact *dbmodel.Contacts) *model.ContactResponse {
	return &model.ContactResponse{
		ID:        contact.ID,
		FirstName: contact.FirstName,
		LastName:  stringValueOrEmpty(contact.LastName),
		Email:     stringValueOrEmpty(contact.Email),
		Phone:     stringValueOrEmpty(contact.Phone),
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}
}
