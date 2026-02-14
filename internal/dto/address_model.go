package dto

type AddressResponse struct {
	ID         string `json:"id"`
	Street     string `json:"street"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
	CreatedAt  int64  `json:"createdAt"`
	UpdatedAt  int64  `json:"updatedAt"`
}

type ListAddressRequest struct {
	UserId    string `json:"-" validate:"required,max=36,uuid"`
	ContactId string `json:"-" validate:"required,max=36,uuid"`
}

type CreateAddressRequest struct {
	UserId     string `json:"-" validate:"required,max=36,uuid"`
	ContactId  string `json:"-" validate:"required,max=36,uuid"`
	Street     string `json:"street" validate:"max=255"`
	City       string `json:"city" validate:"max=255"`
	Province   string `json:"province" validate:"max=255"`
	PostalCode string `json:"postalCode" validate:"max=10"`
	Country    string `json:"country" validate:"max=100"`
}

type UpdateAddressRequest struct {
	UserId     string `json:"-" validate:"required,max=36,uuid"`
	ContactId  string `json:"-" validate:"required,max=36,uuid"`
	ID         string `json:"-" validate:"required,max=36,uuid"`
	Street     string `json:"street" validate:"max=255"`
	City       string `json:"city" validate:"max=255"`
	Province   string `json:"province" validate:"max=255"`
	PostalCode string `json:"postalCode" validate:"max=10"`
	Country    string `json:"country" validate:"max=100"`
}

type GetAddressRequest struct {
	UserId    string `json:"-" validate:"required,max=36,uuid"`
	ContactId string `json:"-" validate:"required,max=36,uuid"`
	ID        string `json:"-" validate:"required,max=36,uuid"`
}

type DeleteAddressRequest struct {
	UserId    string `json:"-" validate:"required,max=36,uuid"`
	ContactId string `json:"-" validate:"required,max=36,uuid"`
	ID        string `json:"-" validate:"required,max=36,uuid"`
}
