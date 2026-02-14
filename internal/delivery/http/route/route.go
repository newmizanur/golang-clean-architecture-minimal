package route

import (
	"golang-clean-architecture/internal/delivery/http"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	App               *echo.Echo
	UserController    *http.UserController
	ContactController *http.ContactController
	AddressController *http.AddressController
	AuthMiddleware    echo.MiddlewareFunc
	ItemController    *http.ItemController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.POST("/api/users", c.UserController.Register)
	c.App.POST("/api/users/login", c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute() {
	auth := c.App.Group("", c.AuthMiddleware)
	auth.DELETE("/api/users", c.UserController.Logout)
	auth.PATCH("/api/users/current", c.UserController.Update)
	auth.GET("/api/users/current", c.UserController.Current)

	auth.GET("/api/contacts", c.ContactController.List)
	auth.POST("/api/contacts", c.ContactController.Create)
	auth.PUT("/api/contacts/:contactId", c.ContactController.Update)
	auth.GET("/api/contacts/:contactId", c.ContactController.Get)
	auth.DELETE("/api/contacts/:contactId", c.ContactController.Delete)

	auth.GET("/api/contacts/:contactId/addresses", c.AddressController.List)
	auth.POST("/api/contacts/:contactId/addresses", c.AddressController.Create)
	auth.PUT("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Update)
	auth.GET("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Get)
	auth.DELETE("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Delete)

	auth.GET("/api/items", c.ItemController.List)
	auth.POST("/api/items", c.ItemController.Create)
	auth.GET("/api/items/:itemId", c.ItemController.Get)
	auth.PUT("/api/items/:itemId", c.ItemController.Update)
	auth.DELETE("/api/items/:itemId", c.ItemController.Delete)
}
