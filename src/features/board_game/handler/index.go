package handler

import (
	"main/features/board_game/repository"
	"main/features/board_game/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewBoardGameHandler(c *echo.Echo) {
	NewFindItSoloPlayBoardGameHandler(c, usecase.NewFindItSoloPlayBoardGameUseCase(repository.NewFindItSoloPlayBoardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
