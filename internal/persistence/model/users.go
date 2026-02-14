package model

import "github.com/uptrace/bun"

type Users struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        string `bun:",pk,column:id"`
	Name      string `bun:"column:name"`
	Password  string `bun:"column:password"`
	CreatedAt int64  `bun:"column:created_at"`
	UpdatedAt int64  `bun:"column:updated_at"`
}
