package _interface

import (
	"context"
	"main/utils/db/mysql"
)

type IMessageChatRepository interface {
	FindOneChat(ctx context.Context, secret string) (*mysql.Chat, error)
}

type IAuthChatRepository interface {
	InsertOneChat(ctx context.Context, chatDTO *mysql.Chat) error
	FindOneUserInfo(ctx context.Context, userID uint) (*mysql.Users, error)
}
