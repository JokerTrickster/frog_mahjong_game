package features

import (
	authHandler "main/features/auth/handler"
	chatHandler "main/features/chat/handler"
	gameHandler "main/features/game/handler"
	profileHandler "main/features/profiles/handler"
	roomsHandler "main/features/rooms/handler"
	userHandler "main/features/users/handler"
	"main/features/ws"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitHandler(e *echo.Echo) error {
	//elb 헬스체크용
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	//인증 핸들러 초기화
	authHandler.NewAuthHandler(e)
	roomsHandler.NewRoomsHandler(e)
	gameHandler.NewGameHandler(e)
	userHandler.NewUsersHandler(e)
	chatHandler.NewChatHandler(e)
	profileHandler.NewProfilesHandler(e)
	//websocket 초기화
	ws.NewWebsocketHandler(e)
	go ws.WSHandleMessages()

	return nil
}
