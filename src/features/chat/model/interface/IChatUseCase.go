package _interface

import (
	"context"
	"main/utils/db/mysql"
)

type IMessageChatUseCase interface {
	Message(c context.Context, secret string) (*mysql.Chat, error)
}

type IAuthChatUseCase interface {
	Auth(c context.Context, userID uint) (string, error)
}
