package http

import (
	"golang-clean-architecture/internal/apperror"
	"golang-clean-architecture/internal/delivery/http/middleware"
	httpresponse "golang-clean-architecture/internal/delivery/http/response"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ContactController struct {
	UseCase *usecase.ContactUseCase
	Log     *logrus.Logger
}

func NewContactController(useCase *usecase.ContactUseCase, log *logrus.Logger) *ContactController {
	return &ContactController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *ContactController) Create(ctx echo.Context) error {
	auth := middleware.GetUser(ctx)

	request := new(model.CreateContactRequest)
	if err := ctx.Bind(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return httpresponse.NewErrorBuilder(apperror.ContactErrors.InvalidRequest).Send(ctx)
	}
	request.UserId = auth.ID

	response, err := c.UseCase.Create(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating contact")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(response).Send(ctx)
}

func (c *ContactController) List(ctx echo.Context) error {
	auth := middleware.GetUser(ctx)

	request := &model.SearchContactRequest{
		UserId: auth.ID,
		Name:   ctx.QueryParam("name"),
		Email:  ctx.QueryParam("email"),
		Phone:  ctx.QueryParam("phone"),
		Page:   intParam(ctx, "page", 1),
		Size:   intParam(ctx, "size", 10),
	}

	responses, total, err := c.UseCase.Search(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching contact")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return httpresponse.SuccessBuilder(responses).WithPaging(paging).Send(ctx)
}

func (c *ContactController) Get(ctx echo.Context) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetContactRequest{
		UserId: auth.ID,
		ID:     ctx.Param("contactId"),
	}

	response, err := c.UseCase.Get(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(response).Send(ctx)
}

func (c *ContactController) Update(ctx echo.Context) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateContactRequest)
	if err := ctx.Bind(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return httpresponse.NewErrorBuilder(apperror.ContactErrors.InvalidRequest).Send(ctx)
	}

	request.UserId = auth.ID
	request.ID = ctx.Param("contactId")

	response, err := c.UseCase.Update(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(response).Send(ctx)
}

func (c *ContactController) Delete(ctx echo.Context) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Param("contactId")

	request := &model.DeleteContactRequest{
		UserId: auth.ID,
		ID:     contactId,
	}

	if err := c.UseCase.Delete(ctx.Request().Context(), request); err != nil {
		c.Log.WithError(err).Error("error deleting contact")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(true).Send(ctx)
}

func intParam(ctx echo.Context, key string, fallback int) int {
	value, err := strconv.Atoi(ctx.QueryParam(key))
	if err != nil {
		return fallback
	}
	return value
}
