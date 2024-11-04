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
	NewWinRequestGameHandler(c, usecase.NewWinRequestGameUseCase(repository.NewWinRequestGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewResultGameHandler(c, usecase.NewResultGameUseCase(repository.NewResultGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewReportGameHandler(c, usecase.NewReportGameUseCase(repository.NewReportGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewMetaGameHandler(c, usecase.NewMetaGameUseCase(repository.NewMetaGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewDeckCardGameHandler(c, usecase.NewDeckCardGameUseCase(repository.NewDeckCardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewListMissionGameHandler(c, usecase.NewListMissionGameUseCase(repository.NewListMissionGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewV2ReportGameHandler(c, usecase.NewV2ReportGameUseCase(repository.NewV2ReportGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewCreateMissionGameHandler(c, usecase.NewCreateMissionGameUseCase(repository.NewCreateMissionGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewListCardGameHandler(c, usecase.NewListCardGameUseCase(repository.NewListCardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
