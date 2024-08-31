package _interface

import (
	"context"
	"main/features/chat/model/request"
	"main/features/chat/model/response"
	"main/utils/db/mysql"
)

type IMessageChatUseCase interface {
	Message(c context.Context, secret string) (*mysql.Chats, error)
}

type IAuthChatUseCase interface {
	Auth(c context.Context, userID uint) (string, error)
}

type IHistoryChatUseCase interface {
	History(c context.Context, req *request.ReqHistory) (response.ResHistoryChat, error)
}
