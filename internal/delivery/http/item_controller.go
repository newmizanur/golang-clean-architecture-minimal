package http

import (
	"golang-clean-architecture/internal/apperror"
	"golang-clean-architecture/internal/delivery/http/middleware"
	httpresponse "golang-clean-architecture/internal/delivery/http/response"
	"golang-clean-architecture/internal/dto"
	"golang-clean-architecture/internal/usecase"
	"math"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ItemController struct {
	ItemUseCase *usecase.ItemUseCase
	Log         *logrus.Logger
}

func NewItemController(useCase *usecase.ItemUseCase, log *logrus.Logger) *ItemController {
	return &ItemController{
		ItemUseCase: useCase,
		Log:         log,
	}
}

func (c *ItemController) Create(ctx echo.Context) error {
	_ = middleware.GetUser(ctx)

	request := new(dto.CreateItemRequest)
	if err := ctx.Bind(&request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return httpresponse.NewErrorBuilder(apperror.ItemErrors.InvalidRequest).Send(ctx)
	}

	response, err := c.ItemUseCase.Create(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error on item create")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(response).Send(ctx)
}

func (c *ItemController) List(ctx echo.Context) error {
	_ = middleware.GetUser(ctx)

	request := &dto.SearchItemRequest{
		Name: ctx.QueryParam("name"),
		SKU:  ctx.QueryParam("sku"),
		Sort: ctx.QueryParam("sort"),
		Page: intParam(ctx, "page", 1),
		Size: intParam(ctx, "size", 10),
	}

	response, total, err := c.ItemUseCase.Search(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error on search")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	paging := &dto.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return httpresponse.SuccessBuilder(response).WithPaging(paging).Send(ctx)
}
