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
	FindById(ctx context.Context, tx bun.IDB, id string) (*m.User, error)
	Create(ctx context.Context, tx bun.IDB, user *m.User) error
	Update(ctx context.Context, tx bun.IDB, user *m.User) error
	Delete(ctx context.Context, tx bun.IDB, user *m.User) error
}

// ContactRepositoryPort defines the contract for contact persistence operations.
type ContactRepositoryPort interface {
	FindByIdAndUserId(ctx context.Context, tx bun.IDB, id string, userId string) (*m.Contact, error)
	Search(ctx context.Context, tx bun.IDB, request *dto.SearchContactRequest) ([]m.Contact, int64, error)
	Create(ctx context.Context, tx bun.IDB, contact *m.Contact) error
	Update(ctx context.Context, tx bun.IDB, contact *m.Contact) error
	Delete(ctx context.Context, tx bun.IDB, contact *m.Contact) error
}

// AddressRepositoryPort defines the contract for address persistence operations.
type AddressRepositoryPort interface {
	FindByIdAndContactId(ctx context.Context, tx bun.IDB, id string, contactId string) (*m.Address, error)
	FindAllByContactId(ctx context.Context, tx bun.IDB, contactId string) ([]m.Address, error)
	Create(ctx context.Context, tx bun.IDB, address *m.Address) error
	Update(ctx context.Context, tx bun.IDB, address *m.Address) error
	Delete(ctx context.Context, tx bun.IDB, address *m.Address) error
}

// ItemRepositoryPort defines the contract for item persistence operations.
type ItemRepositoryPort interface {
	FindById(ctx context.Context, tx bun.IDB, id int64) (*m.Item, error)
	Search(ctx context.Context, tx bun.IDB, search *dto.SearchItemRequest) ([]m.Item, int64, error)
	Create(ctx context.Context, tx bun.IDB, item *m.Item) (int64, error)
	Update(ctx context.Context, tx bun.IDB, item *m.Item) error
	Delete(ctx context.Context, tx bun.IDB, item *m.Item) error
}
