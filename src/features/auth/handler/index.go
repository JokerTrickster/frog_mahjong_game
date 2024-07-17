package handler

import (
	"main/features/auth/repository"
	"main/features/auth/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewAuthHandler(c *echo.Echo) {
	NewSignupAuthHandler(c, usecase.NewSignupAuthUseCase(repository.NewSignupAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewSigninAuthHandler(c, usecase.NewSigninAuthUseCase(repository.NewSigninAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewLogoutAuthHandler(c, usecase.NewLogoutAuthUseCase(repository.NewLogoutAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewReissueAuthHandler(c, usecase.NewReissueAuthUseCase(repository.NewReissueAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewGoogleOauthAuthHandler(c, usecase.NewGoogleOauthAuthUseCase(repository.NewGoogleOauthAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewGoogleOauthCallbackAuthHandler(c, usecase.NewGoogleOauthCallbackAuthUseCase(repository.NewGoogleOauthCallbackAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewRequestPasswordAuthHandler(c, usecase.NewRequestPasswordAuthUseCase(repository.NewRequestPasswordAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewValidatePasswordAuthHandler(c, usecase.NewValidatePasswordAuthUseCase(repository.NewValidatePasswordAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
