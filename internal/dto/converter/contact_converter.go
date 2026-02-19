package converter

import (
	"golang-clean-architecture/ent"
	"golang-clean-architecture/internal/dto"
)

func ContactToResponse(contact *ent.Contact) *dto.ContactResponse {
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
