package _interface

import (
	"context"
	"main/utils/db/mysql"
)

type IMessageChatUseCase interface {
	Message(c context.Context, secret string) (*mysql.Chats, error)
}

type IAuthChatUseCase interface {
	Auth(c context.Context, userID uint) (string, error)
}