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
	NewDiscardGameHandler(c, usecase.NewDiscardGameUseCase(repository.NewDiscardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewNextTurnGameHandler(c, usecase.NewNextTurnGameUseCase(repository.NewNextTurnGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewLoanGameHandler(c, usecase.NewLoanGameUseCase(repository.NewLoanGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewScoreCalculateGameHandler(c, usecase.NewScoreCalculateGameUseCase(repository.NewScoreCalculateGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
