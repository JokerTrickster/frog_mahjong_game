package handler

import (
	"main/features/board_game/repository"
	"main/features/board_game/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewBoardGameHandler(c *echo.Echo) {
	NewFindItSoloPlayBoardGameHandler(c, usecase.NewFindItSoloPlayBoardGameUseCase(repository.NewFindItSoloPlayBoardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewFindItRankBoardGameHandler(c, usecase.NewFindItRankBoardGameUseCase(repository.NewFindItRankBoardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewFindItCoinBoardGameHandler(c, usecase.NewFindItCoinBoardGameUseCase(repository.NewFindItCoinBoardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewFindItPasswordCheckBoardGameHandler(c, usecase.NewFindItPasswordCheckBoardGameUseCase(repository.NewFindItPasswordCheckBoardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewSlimeWarGetsCardBoardGameHandler(c, usecase.NewSlimeWarGetsCardBoardGameUseCase(repository.NewSlimeWarGetsCardBoardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewSlimeWarResultBoardGameHandler(c, usecase.NewSlimeWarResultBoardGameUseCase(repository.NewSlimeWarResultBoardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewSlimeWarRankBoardGameHandler(c, usecase.NewSlimeWarRankBoardGameUseCase(repository.NewSlimeWarRankBoardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewSequenceResultBoardGameHandler(c, usecase.NewSequenceResultBoardGameUseCase(repository.NewSequenceResultBoardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewSequenceRankBoardGameHandler(c, usecase.NewSequenceRankBoardGameUseCase(repository.NewSequenceRankBoardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewGameOverBoardGameHandler(c, usecase.NewGameOverBoardGameUseCase(repository.NewGameOverBoardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewFrogCardListBoardGameHandler(c, usecase.NewFrogCardListBoardGameUseCase(repository.NewFrogCardListBoardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
