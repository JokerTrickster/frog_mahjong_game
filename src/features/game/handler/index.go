package handler

import (
	"main/features/game/repository"
	"main/features/game/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewGameHandler(c *echo.Echo) {
	NewStartGameHandler(c, usecase.NewStartGameUseCase(repository.NewStartGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
