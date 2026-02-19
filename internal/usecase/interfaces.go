package usecase

import (
	"context"

	"golang-clean-architecture/ent"
	"golang-clean-architecture/internal/dto"
)

// UserRepositoryPort defines the contract for user persistence operations.
type UserRepositoryPort interface {
	CountById(ctx context.Context, client *ent.Client, id string) (int64, error)
	FindById(ctx context.Context, client *ent.Client, id string) (*ent.User, error)
	Create(ctx context.Context, client *ent.Client, user *ent.User) error
	Update(ctx context.Context, client *ent.Client, user *ent.User) error
	Delete(ctx context.Context, client *ent.Client, id string) error
}

// ContactRepositoryPort defines the contract for contact persistence operations.
type ContactRepositoryPort interface {
	FindByIdAndUserId(ctx context.Context, client *ent.Client, id string, userId string) (*ent.Contact, error)
	Search(ctx context.Context, client *ent.Client, request *dto.SearchContactRequest) ([]*ent.Contact, int64, error)
	Create(ctx context.Context, client *ent.Client, contact *ent.Contact) error
	Update(ctx context.Context, client *ent.Client, contact *ent.Contact) error
	Delete(ctx context.Context, client *ent.Client, id string) error
}

// AddressRepositoryPort defines the contract for address persistence operations.
type AddressRepositoryPort interface {
	FindByIdAndContactId(ctx context.Context, client *ent.Client, id string, contactId string) (*ent.Address, error)
	FindAllByContactId(ctx context.Context, client *ent.Client, contactId string) ([]*ent.Address, error)
	Create(ctx context.Context, client *ent.Client, address *ent.Address) error
	Update(ctx context.Context, client *ent.Client, address *ent.Address) error
	Delete(ctx context.Context, client *ent.Client, id string) error
}

// ItemRepositoryPort defines the contract for item persistence operations.
type ItemRepositoryPort interface {
	FindById(ctx context.Context, client *ent.Client, id int) (*ent.Item, error)
	Search(ctx context.Context, client *ent.Client, search *dto.SearchItemRequest) ([]*ent.Item, int64, error)
	Create(ctx context.Context, client *ent.Client, item *ent.Item) (*ent.Item, error)
	Update(ctx context.Context, client *ent.Client, item *ent.Item) error
	Delete(ctx context.Context, client *ent.Client, id int) error
}
