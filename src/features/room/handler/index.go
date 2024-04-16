package handler

import (
	"main/features/room/repository"
	"main/features/room/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewRoomHandler(c *echo.Echo) {
	NewCreateRoomHandler(c, usecase.NewCreateRoomUseCase(repository.NewCreateRoomRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewJoinRoomHandler(c, usecase.NewJoinRoomUseCase(repository.NewJoinRoomRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewOutRoomHandler(c, usecase.NewOutRoomUseCase(repository.NewOutRoomRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewReadyRoomHandler(c, usecase.NewReadyRoomUseCase(repository.NewReadyRoomRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewListRoomHandler(c, usecase.NewListRoomUseCase(repository.NewListRoomRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewLogoutRoomHandler(c, usecase.NewLogoutRoomUseCase(repository.NewLogoutRoomRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
