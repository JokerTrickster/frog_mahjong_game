package features

import (
	authHandler "main/features/auth/handler"
	boardGameHandler "main/features/board_game/handler"
	chatHandler "main/features/chat/handler"
	"main/features/find_it"
	frog "main/features/frog"
	gameHandler "main/features/game/handler"
	gameAuthHandler "main/features/game_auth/handler"
	gameProfileHandler "main/features/game_profiles/handler"
	gameUserHandler "main/features/game_users/handler"
	profileHandler "main/features/profiles/handler"
	roomsHandler "main/features/rooms/handler"
	"main/features/sequence"
	"main/features/slime_war"
	slimeWar "main/features/slime_war"
	userHandler "main/features/users/handler"
	"main/features/v2ws"
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
	gameAuthHandler.NewGameAuthHandler(e)
	gameUserHandler.NewGameUsersHandler(e)
	gameProfileHandler.NewGameProfilesHandler(e)
	boardGameHandler.NewBoardGameHandler(e)
	//websocket 초기화
	v2ws.NewV2WebsocketHandler(e)
	find_it.NewFindItWebsocketHandler(e)
	slimeWar.NewSlimeWarWebsocketHandler(e)
	sequence.NewSequenceWebsocketHandler(e)
	frog.NewWebsocketHandler(e)

	go v2ws.WSHandleMessages("wingspan")

	go find_it.WSHandleMessages("find-it")

	go slime_war.WSHandleMessages("slime-war")

	go sequence.WSHandleMessages("sequence")

	go frog.WSHandleMessages("frog")

	return nil
}
