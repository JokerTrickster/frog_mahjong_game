package _interface

import (
	"context"
	"main/utils/db/mysql"
)

type IMessageChatRepository interface {
	FindOneChat(ctx context.Context, secret string) (*mysql.Chats, error)
}

type IAuthChatRepository interface {
	InsertOneChat(ctx context.Context, chatDTO *mysql.Chats) error
	FindOneUserInfo(ctx context.Context, userID uint) (*mysql.Users, error)
}
