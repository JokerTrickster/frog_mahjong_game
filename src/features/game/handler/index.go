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
	NewWinRequestGameHandler(c, usecase.NewWinRequestGameUseCase(repository.NewWinRequestGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewResultGameHandler(c, usecase.NewResultGameUseCase(repository.NewResultGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewReportGameHandler(c, usecase.NewReportGameUseCase(repository.NewReportGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewMetaGameHandler(c, usecase.NewMetaGameUseCase(repository.NewMetaGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewDeckCardGameHandler(c, usecase.NewDeckCardGameUseCase(repository.NewDeckCardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewListMissionGameHandler(c, usecase.NewListMissionGameUseCase(repository.NewListMissionGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewV2ReportGameHandler(c, usecase.NewV2ReportGameUseCase(repository.NewV2ReportGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewCreateMissionGameHandler(c, usecase.NewCreateMissionGameUseCase(repository.NewCreateMissionGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewV2ListCardGameHandler(c, usecase.NewV2ListCardGameUseCase(repository.NewV2ListCardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewV2DeckCardGameHandler(c, usecase.NewV2DeckCardGameUseCase(repository.NewV2DeckCardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewV2ResultGameHandler(c, usecase.NewV2ResultGameUseCase(repository.NewV2ResultGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewSaveCardInfoGameHandler(c, usecase.NewSaveCardInfoGameUseCase(repository.NewSaveCardInfoGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewSaveCardImageGameHandler(c, usecase.NewSaveCardImageGameUseCase(repository.NewSaveCardImageGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewUpdateCardGameHandler(c, usecase.NewUpdateCardGameUseCase(repository.NewUpdateCardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewReportImageUploadGameHandler(c, usecase.NewReportImageUploadGameUseCase(repository.NewReportImageUploadGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewListCardGameHandler(c, usecase.NewListCardGameUseCase(repository.NewListCardGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewV2DrawResultGameHandler(c, usecase.NewV2DrawResultGameUseCase(repository.NewV2DrawResultGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))

	//find-it
	NewFindItResultGameHandler(c, usecase.NewFindItResultGameUseCase(repository.NewFindItResultGameRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
