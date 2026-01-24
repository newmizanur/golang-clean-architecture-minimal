package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Item struct {
	bun.BaseModel `bun:"table:items,alias:i"`

	ID        int64      `bun:",pk,autoincrement,column:id"`
	Name      string     `bun:"column:name"`
	Sku       string     `bun:"column:sku"`
	Currency  string     `bun:"column:currency"`
	Stock     int32      `bun:"column:stock"`
	CreatedAt *time.Time `bun:"column:created_at"`
	UpdatedAt *time.Time `bun:"column:updated_at"`
}
