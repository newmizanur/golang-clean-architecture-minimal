package model

type ContactResponse struct {
	ID        string            `json:"id"`
	FirstName string            `json:"firstName"`
	LastName  string            `json:"lastName"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	CreatedAt int64             `json:"createdAt"`
	UpdatedAt int64             `json:"updatedAt"`
	Addresses []AddressResponse `json:"addresses,omitempty"`
}

type CreateContactRequest struct {
	UserId    string `json:"-" validate:"required,max=36,uuid"`
	FirstName string `json:"firstName" validate:"required,max=100"`
	LastName  string `json:"lastName" validate:"max=100"`
	Email     string `json:"email" validate:"max=200,email"`
	Phone     string `json:"phone" validate:"max=20"`
}

type UpdateContactRequest struct {
	UserId    string `json:"-" validate:"required,max=36,uuid"`
	ID        string `json:"-" validate:"required,max=36,uuid"`
	FirstName string `json:"firstName" validate:"required,max=100"`
	LastName  string `json:"lastName" validate:"max=100"`
	Email     string `json:"email" validate:"max=200,email"`
	Phone     string `json:"phone" validate:"max=20"`
}

type SearchContactRequest struct {
	UserId string `json:"-" validate:"required,max=36,uuid"`
	Name   string `json:"name" validate:"max=100"`
	Email  string `json:"email" validate:"max=200"`
	Phone  string `json:"phone" validate:"max=20"`
	Page   int    `json:"page" validate:"min=1"`
	Size   int    `json:"size" validate:"min=1,max=100"`
}

type GetContactRequest struct {
	UserId string `json:"-" validate:"required,max=36,uuid"`
	ID     string `json:"-" validate:"required,max=36,uuid"`
}

type DeleteContactRequest struct {
	UserId string `json:"-" validate:"required,max=36,uuid"`
	ID     string `json:"-" validate:"required,max=36,uuid"`
}
