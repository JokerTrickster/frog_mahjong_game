package handler

import (
	"main/features/chat/repository"
	"main/features/chat/usecase"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
)

func NewChatHandler(c *echo.Echo) {
	NewMessageChatHandler(c, usecase.NewMessageChatUseCase(repository.NewMessageChatRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewAuthChatHandler(c, usecase.NewAuthChatUseCase(repository.NewAuthChatRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
	NewHistoryChatHandler(c, usecase.NewHistoryChatUseCase(repository.NewHistoryChatRepository(mysql.GormMysqlDB), mysql.DBTimeOut))
}
