package handler

import (
	"main/features/users/repository"
	"main/features/users/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewUsersHandler(c *echo.Echo) {
	NewGetUsersHandler(c, usecase.NewGetUsersUseCase(repository.NewGetUsersRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewListUsersHandler(c, usecase.NewListUsersUseCase(repository.NewListUsersRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewUpdateUsersHandler(c, usecase.NewUpdateUsersUseCase(repository.NewUpdateUsersRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewDeleteUsersHandler(c, usecase.NewDeleteUsersUseCase(repository.NewDeleteUsersRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
