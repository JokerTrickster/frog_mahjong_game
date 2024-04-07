package handler

import (
	"main/features/game/repository"
	"main/features/game/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewGameHandler(c *echo.Echo) {
	NewStartGameHandler(c, usecase.NewStartGameUseCase(repository.NewStartGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewDoraGameHandler(c, usecase.NewDoraGameUseCase(repository.NewDoraGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewOwnershipGameHandler(c, usecase.NewOwnershipGameUseCase(repository.NewOwnershipGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
