package config

import (
	"net/http"

	"golang-clean-architecture/internal/apperror"
	"golang-clean-architecture/internal/delivery/http/response"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func NewEcho(config *viper.Viper) *echo.Echo {
	_ = config
	app := echo.New()
	app.HideBanner = true
	app.HTTPErrorHandler = NewErrorHandler()

	return app
}

func NewErrorHandler() echo.HTTPErrorHandler {
	return func(err error, ctx echo.Context) {
		if ctx.Response().Committed {
			return
		}

		if appErr, ok := err.(*apperror.AppError); ok {
			if sendErr := response.NewErrorBuilder(appErr).Send(ctx); sendErr != nil {
				ctx.Logger().Error(sendErr)
			}
			return
		}

		if httpErr, ok := err.(*echo.HTTPError); ok {
			message := "Internal Server Error"
			switch value := httpErr.Message.(type) {
			case string:
				message = value
			case error:
				message = value.Error()
			default:
				message = "Error"
			}

			fallback := apperror.NewAppError(httpErr.Code, httpErr.Code, message)
			if sendErr := response.NewErrorBuilder(fallback).Send(ctx); sendErr != nil {
				ctx.Logger().Error(sendErr)
			}
			return
		}

		fallback := apperror.NewAppError(http.StatusInternalServerError, http.StatusInternalServerError, "internal server error")
		if sendErr := response.NewErrorBuilder(fallback).Send(ctx); sendErr != nil {
			ctx.Logger().Error(sendErr)
		}
	}
}
