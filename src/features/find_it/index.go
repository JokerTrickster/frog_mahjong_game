package find_it

import "github.com/labstack/echo/v4"

func NewFindItWebsocketHandler(e *echo.Echo) {

	// e.GET("/v0.1/rooms/join/ws", join)
	e.GET("/find-it/v0.1/rooms/match/ws", match)
	e.GET("/find-it/v0.1/rooms/play/together/ws", playTogether)
	e.GET("/find-it/v0.1/rooms/join/ws", joinPlay)
}
