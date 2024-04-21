package handler

import (
	"main/features/rooms/repository"
	"main/features/rooms/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewRoomsHandler(c *echo.Echo) {
	NewCreateRoomsHandler(c, usecase.NewCreateRoomsUseCase(repository.NewCreateRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewJoinRoomsHandler(c, usecase.NewJoinRoomsUseCase(repository.NewJoinRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewOutRoomsHandler(c, usecase.NewOutRoomsUseCase(repository.NewOutRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewReadyRoomsHandler(c, usecase.NewReadyRoomsUseCase(repository.NewReadyRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewListRoomsHandler(c, usecase.NewListRoomsUseCase(repository.NewListRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewUserListRoomsHandler(c, usecase.NewUserListRoomsUseCase(repository.NewUserListRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
