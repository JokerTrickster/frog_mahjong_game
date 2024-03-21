package handler

import (
	"main/features/room/repository"
	"main/features/room/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewRoomHandler(c *echo.Echo) {
	NewCreateRoomHandler(c, usecase.NewCreateRoomUseCase(repository.NewCreateRoomRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
