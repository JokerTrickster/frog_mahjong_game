package sequence

import "github.com/labstack/echo/v4"

func NewSequenceWebsocketHandler(e *echo.Echo) {

	// e.GET("/v0.1/rooms/join/ws", join)
	e.GET("/sequence/v0.1/rooms/match/ws", match)
	e.GET("/sequence/v0.1/rooms/play/together/ws", playTogether)
	e.GET("/sequence/v0.1/rooms/play/join/ws", joinPlay)
}
