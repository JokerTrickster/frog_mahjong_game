package handler

import (
	"main/features/profiles/repository"
	"main/features/profiles/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewProfilesHandler(c *echo.Echo) {
	NewListProfilesHandler(c, usecase.NewListProfilesUseCase(repository.NewListProfilesRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewUploadProfilesHandler(c, usecase.NewUploadProfilesUseCase(repository.NewUploadProfilesRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
