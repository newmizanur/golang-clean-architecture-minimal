package http

import (
	"golang-clean-architecture/internal/apperror"
	"golang-clean-architecture/internal/delivery/http/middleware"
	httpresponse "golang-clean-architecture/internal/delivery/http/response"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AddressController struct {
	UseCase *usecase.AddressUseCase
	Log     *logrus.Logger
}

func NewAddressController(useCase *usecase.AddressUseCase, log *logrus.Logger) *AddressController {
	return &AddressController{
		Log:     log,
		UseCase: useCase,
	}
}

func (c *AddressController) Create(ctx echo.Context) error {
	auth := middleware.GetUser(ctx)

	request := new(model.CreateAddressRequest)
	if err := ctx.Bind(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return httpresponse.NewErrorBuilder(apperror.AddressErrors.InvalidRequest).Send(ctx)
	}

	request.UserId = auth.ID
	request.ContactId = ctx.Param("contactId")

	response, err := c.UseCase.Create(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create address")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(response).Send(ctx)
}

func (c *AddressController) List(ctx echo.Context) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Param("contactId")

	request := &model.ListAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
	}

	responses, err := c.UseCase.List(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list addresses")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(responses).Send(ctx)
}

func (c *AddressController) Get(ctx echo.Context) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Param("contactId")
	addressId := ctx.Param("addressId")

	request := &model.GetAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
		ID:        addressId,
	}

	response, err := c.UseCase.Get(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to get address")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(response).Send(ctx)
}

func (c *AddressController) Update(ctx echo.Context) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateAddressRequest)
	if err := ctx.Bind(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return httpresponse.NewErrorBuilder(apperror.AddressErrors.InvalidRequest).Send(ctx)
	}

	request.UserId = auth.ID
	request.ContactId = ctx.Param("contactId")
	request.ID = ctx.Param("addressId")

	response, err := c.UseCase.Update(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to update address")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(response).Send(ctx)
}

func (c *AddressController) Delete(ctx echo.Context) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Param("contactId")
	addressId := ctx.Param("addressId")

	request := &model.DeleteAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
		ID:        addressId,
	}

	if err := c.UseCase.Delete(ctx.Request().Context(), request); err != nil {
		c.Log.WithError(err).Error("failed to delete address")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(true).Send(ctx)
}
