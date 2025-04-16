package slime_war

import "github.com/labstack/echo/v4"

func NewFindItWebsocketHandler(e *echo.Echo) {

	// e.GET("/v0.1/rooms/join/ws", join)
	e.GET("/slime-war/v0.1/rooms/match/ws", match)
	e.GET("/slime-war/v0.1/rooms/play/together/ws", playTogether)
	e.GET("/slime-war/v0.1/rooms/play/join/ws", joinPlay)
}
