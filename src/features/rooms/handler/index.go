package handler

import (
	"main/features/rooms/repository"
	"main/features/rooms/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewRoomsHandler(c *echo.Echo) {
	NewCreateRoomsHandler(c, usecase.NewCreateRoomsUseCase(repository.NewCreateRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewJoinPlayRoomsHandler(c, usecase.NewJoinPlayRoomsUseCase(repository.NewJoinPlayRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewOutRoomsHandler(c, usecase.NewOutRoomsUseCase(repository.NewOutRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewReadyRoomsHandler(c, usecase.NewReadyRoomsUseCase(repository.NewReadyRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewListRoomsHandler(c, usecase.NewListRoomsUseCase(repository.NewListRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewUserListRoomsHandler(c, usecase.NewUserListRoomsUseCase(repository.NewUserListRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewV02CreateRoomsHandler(c, usecase.NewV02CreateRoomsUseCase(repository.NewV02CreateRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewV02JoinRoomsHandler(c, usecase.NewV02JoinRoomsUseCase(repository.NewV02JoinRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewMetaRoomsHandler(c, usecase.NewMetaRoomsUseCase(repository.NewMetaRoomsRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
