package model

import "github.com/uptrace/bun"

type Address struct {
	bun.BaseModel `bun:"table:addresses,alias:a"`

	ID         string  `bun:",pk,column:id"`
	ContactID  string  `bun:"column:contact_id"`
	Street     *string `bun:"column:street"`
	City       *string `bun:"column:city"`
	Province   *string `bun:"column:province"`
	PostalCode *string `bun:"column:postal_code"`
	Country    *string `bun:"column:country"`
	CreatedAt  int64   `bun:"column:created_at"`
	UpdatedAt  int64   `bun:"column:updated_at"`
}
