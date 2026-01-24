package response

import (
	"net/http"

	"golang-clean-architecture/internal/model"

	"github.com/labstack/echo/v4"
)

type SuccessResponse[T any] struct {
	Data   T
	Paging *model.PageMetadata
}

func SuccessBuilder[T any](data T) *SuccessResponse[T] {
	return &SuccessResponse[T]{Data: data}
}

func (r *SuccessResponse[T]) WithPaging(paging *model.PageMetadata) *SuccessResponse[T] {
	r.Paging = paging
	return r
}

func (r *SuccessResponse[T]) Send(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, model.WebResponse[T]{Data: r.Data, Paging: r.Paging})
}
