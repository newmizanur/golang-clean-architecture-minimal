package converter

import (
	"golang-clean-architecture/ent"
	"golang-clean-architecture/internal/dto"
)

func ItemToResponse(item *ent.Item) *dto.CreateItemResponse {
	return &dto.CreateItemResponse{
		ID:        int64(item.ID),
		Name:      item.Name,
		SKU:       item.Sku,
		Currency:  item.Currency,
		Stock:     item.Stock,
		CreatedAt: &item.CreatedAt,
		UpdatedAt: &item.UpdatedAt,
	}
}
