package router

import (
	"github.com/labstack/echo/v4"
	"jougan/handler"
	"net/http"
)

type API struct {
	Echo        *echo.Echo
	PromHandler http.Handler
}

func (api *API) SetupRouter() {

	api.Echo.GET("/", handler.Welcome)
	api.Echo.GET("/metrics", echo.WrapHandler(api.PromHandler))
}
