package usecase

import (
	"context"

	"golang-clean-architecture/internal/dto"
	m "golang-clean-architecture/internal/persistence/model"

	"github.com/uptrace/bun"
)

// UserRepositoryPort defines the contract for user persistence operations.
type UserRepositoryPort interface {
	CountById(ctx context.Context, tx bun.IDB, id string) (int64, error)
	FindById(ctx context.Context, tx bun.IDB, id string) (*m.Users, error)
	Create(ctx context.Context, tx bun.IDB, user *m.Users) error
	Update(ctx context.Context, tx bun.IDB, user *m.Users) error
	Delete(ctx context.Context, tx bun.IDB, user *m.Users) error
}

// ContactRepositoryPort defines the contract for contact persistence operations.
type ContactRepositoryPort interface {
	FindByIdAndUserId(ctx context.Context, tx bun.IDB, id string, userId string) (*m.Contacts, error)
	Search(ctx context.Context, tx bun.IDB, request *dto.SearchContactRequest) ([]m.Contacts, int64, error)
	Create(ctx context.Context, tx bun.IDB, contact *m.Contacts) error
	Update(ctx context.Context, tx bun.IDB, contact *m.Contacts) error
	Delete(ctx context.Context, tx bun.IDB, contact *m.Contacts) error
}

// AddressRepositoryPort defines the contract for address persistence operations.
type AddressRepositoryPort interface {
	FindByIdAndContactId(ctx context.Context, tx bun.IDB, id string, contactId string) (*m.Addresses, error)
	FindAllByContactId(ctx context.Context, tx bun.IDB, contactId string) ([]m.Addresses, error)
	Create(ctx context.Context, tx bun.IDB, address *m.Addresses) error
	Update(ctx context.Context, tx bun.IDB, address *m.Addresses) error
	Delete(ctx context.Context, tx bun.IDB, address *m.Addresses) error
}

// ItemRepositoryPort defines the contract for item persistence operations.
type ItemRepositoryPort interface {
	FindById(ctx context.Context, tx bun.IDB, id int64) (*m.Items, error)
	Search(ctx context.Context, tx bun.IDB, search *dto.SearchItemRequest) ([]m.Items, int64, error)
	Create(ctx context.Context, tx bun.IDB, item *m.Items) (int64, error)
	Update(ctx context.Context, tx bun.IDB, item *m.Items) error
	Delete(ctx context.Context, tx bun.IDB, item *m.Items) error
}
