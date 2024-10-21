package v2ws

import "github.com/labstack/echo/v4"

func NewV2WebsocketHandler(e *echo.Echo) {

	// e.GET("/v0.1/rooms/join/ws", join)
	e.GET("/v2.1/rooms/match/ws", match)
	e.GET("/v2.1/rooms/play/together/ws", playTogether)
	e.GET("/v2.1/rooms/join/play/ws", joinPlay)
}
