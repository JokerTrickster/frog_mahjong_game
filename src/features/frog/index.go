package ws

import "github.com/labstack/echo/v4"

func NewWebsocketHandler(e *echo.Echo) {

	// e.GET("/v0.1/rooms/join/ws", join)
	e.GET("/frog/v0.1/rooms/match/ws", match)
	e.GET("/frog/v0.1/rooms/play/together/ws", playTogether)
	e.GET("/frog/v0.1/rooms/join/play/ws", joinPlay)
}
