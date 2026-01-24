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

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *UserController) Register(ctx echo.Context) error {
	request := new(model.RegisterUserRequest)
	if err := ctx.Bind(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return httpresponse.NewErrorBuilder(apperror.UserErrors.InvalidRequest).Send(ctx)
	}

	response, err := c.UseCase.Create(ctx.Request().Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(response).Send(ctx)
}

func (c *UserController) Login(ctx echo.Context) error {
	request := new(model.LoginUserRequest)
	if err := ctx.Bind(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return httpresponse.NewErrorBuilder(apperror.UserErrors.InvalidRequest).Send(ctx)
	}

	response, err := c.UseCase.Login(ctx.Request().Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(response).Send(ctx)
}

func (c *UserController) Current(ctx echo.Context) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetUserRequest{
		ID: auth.ID,
	}

	response, err := c.UseCase.Current(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get current user")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(response).Send(ctx)
}

func (c *UserController) Logout(ctx echo.Context) error {
	auth := middleware.GetUser(ctx)

	request := &model.LogoutUserRequest{
		ID: auth.ID,
	}

	response, err := c.UseCase.Logout(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to logout user")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(response).Send(ctx)
}

func (c *UserController) Update(ctx echo.Context) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateUserRequest)
	if err := ctx.Bind(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return httpresponse.NewErrorBuilder(apperror.UserErrors.InvalidRequest).Send(ctx)
	}

	request.ID = auth.ID
	response, err := c.UseCase.Update(ctx.Request().Context(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to update user")
		return httpresponse.NewErrorBuilder(err).Send(ctx)
	}

	return httpresponse.SuccessBuilder(response).Send(ctx)
}
