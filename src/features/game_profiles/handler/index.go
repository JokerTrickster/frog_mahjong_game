package handler

import (
	"main/features/game_profiles/repository"
	"main/features/game_profiles/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewGameProfilesHandler(c *echo.Echo) {
	NewListProfilesHandler(c, usecase.NewListProfilesUseCase(repository.NewListProfilesRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewUploadProfilesHandler(c, usecase.NewUploadProfilesUseCase(repository.NewUploadProfilesRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewUpdateProfilesHandler(c, usecase.NewUpdateProfilesUseCase(repository.NewUpdateProfilesRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
