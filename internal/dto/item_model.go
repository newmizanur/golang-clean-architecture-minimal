package dto

import "time"

type CreateItemRequest struct {
	Name     string `json:"name" validate:"required,max=255"`
	SKU      string `json:"sku" validate:"required,max=100"`
	Currency string `json:"currency" validate:"required,max=10"`
	Stock    int32  `json:"stock" validate:"required,min=1"`
}
type CreateItemResponse struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	SKU       string     `json:"sku"`
	Currency  string     `json:"currentcy"`
	Stock     int32      `json:"stock"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type SearchItemRequest struct {
	Name string `json:"name" validate:"required,max=255"`
	SKU  string `json:"sku" validate:"required,max=100"`
	Sort string `json:"sort" validate:"max=50"`
	Page int    `json:"page" validate:"min=1"`
	Size int    `json:"size" validate:"min=1,max=100"`
}

type GetItemRequest struct {
	ID int64 `json:"-" validate:"required,min=1"`
}

type UpdateItemRequest struct {
	ID       int64  `json:"-" validate:"required,min=1"`
	Name     string `json:"name" validate:"max=255"`
	SKU      string `json:"sku" validate:"max=100"`
	Currency string `json:"currency" validate:"max=10"`
	Stock    int32  `json:"stock" validate:"min=0"`
}

type DeleteItemRequest struct {
	ID int64 `json:"-" validate:"required,min=1"`
}
