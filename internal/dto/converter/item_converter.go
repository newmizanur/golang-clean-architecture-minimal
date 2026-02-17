package converter

import (
	"golang-clean-architecture/internal/dto"
	dbmodel "golang-clean-architecture/internal/persistence/model"
)

func ItemToResponse(item *dbmodel.Item) *dto.CreateItemResponse {
	return &dto.CreateItemResponse{
		ID:        int64(item.ID),
		Name:      item.Name,
		SKU:       item.Sku,
		Currency:  item.Currency,
		Stock:     item.Stock,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
