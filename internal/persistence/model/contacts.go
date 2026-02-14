package model

import "github.com/uptrace/bun"

type Contacts struct {
	bun.BaseModel `bun:"table:contacts,alias:c"`

	ID        string  `bun:",pk,column:id"`
	FirstName string  `bun:"column:first_name"`
	LastName  *string `bun:"column:last_name"`
	Email     *string `bun:"column:email"`
	Phone     *string `bun:"column:phone"`
	UserID    string  `bun:"column:user_id"`
	CreatedAt int64   `bun:"column:created_at"`
	UpdatedAt int64   `bun:"column:updated_at"`
}
