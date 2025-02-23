package handler

import (
	"main/features/game_auth/repository"
	"main/features/game_auth/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewGameAuthHandler(c *echo.Echo) {
	NewSignupAuthHandler(c, usecase.NewSignupAuthUseCase(repository.NewSignupAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewSigninAuthHandler(c, usecase.NewSigninAuthUseCase(repository.NewSigninAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewLogoutAuthHandler(c, usecase.NewLogoutAuthUseCase(repository.NewLogoutAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewReissueAuthHandler(c, usecase.NewReissueAuthUseCase(repository.NewReissueAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewGoogleOauthAuthHandler(c, usecase.NewGoogleOauthAuthUseCase(repository.NewGoogleOauthAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewGoogleOauthCallbackAuthHandler(c, usecase.NewGoogleOauthCallbackAuthUseCase(repository.NewGoogleOauthCallbackAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewRequestPasswordAuthHandler(c, usecase.NewRequestPasswordAuthUseCase(repository.NewRequestPasswordAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewValidatePasswordAuthHandler(c, usecase.NewValidatePasswordAuthUseCase(repository.NewValidatePasswordAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewV02GoogleOauthCallbackAuthHandler(c, usecase.NewV02GoogleOauthCallbackAuthUseCase(repository.NewV02GoogleOauthCallbackAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewRequestSignupAuthHandler(c, usecase.NewRequestSignupAuthUseCase(repository.NewRequestSignupAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewFCMTokenAuthHandler(c, usecase.NewFCMTokenAuthUseCase(repository.NewFCMTokenAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewCheckSigninAuthHandler(c, usecase.NewCheckSigninAuthUseCase(repository.NewCheckSigninAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewNameCheckAuthHandler(c, usecase.NewNameCheckAuthUseCase(repository.NewNameCheckAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewValidateSignupAuthHandler(c, usecase.NewValidateSignupAuthUseCase(repository.NewValidateSignupAuthRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
