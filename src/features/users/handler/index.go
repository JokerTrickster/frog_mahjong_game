package handler

import (
	"main/features/users/repository"
	"main/features/users/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewUsersHandler(c *echo.Echo) {
	NewGetUsersHandler(c, usecase.NewGetUsersUseCase(repository.NewGetUsersRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}