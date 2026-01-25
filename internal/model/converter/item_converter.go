package converter

import (
	dbmodel "golang-clean-architecture/internal/entity/db/model"
	"golang-clean-architecture/internal/model"
)

func ItemToResponse(item *dbmodel.Items) *model.CreateItemResponse {
	return &model.CreateItemResponse{
		ID:        int64(item.ID),
		Name:      item.Name,
		SKU:       item.Sku,
		Currency:  item.Currency,
		Stock:     item.Stock,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
