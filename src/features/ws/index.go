package ws

import "github.com/labstack/echo/v4"

func NewWebsocketHandler(e *echo.Echo) {

	e.GET("/v0.1/rooms/join/ws", join)
}