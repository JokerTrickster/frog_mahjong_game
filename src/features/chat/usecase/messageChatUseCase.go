package usecase

import (
	"context"
	_interface "main/features/chat/model/interface"
	"main/utils/db/mysql"
	"time"
)

type MessageChatUseCase struct {
	Repository     _interface.IMessageChatRepository
	ContextTimeout time.Duration
}

func NewMessageChatUseCase(repo _interface.IMessageChatRepository, timeout time.Duration) _interface.IMessageChatUseCase {
	return &MessageChatUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *MessageChatUseCase) Message(c context.Context, secret string) (*mysql.Chat, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	chatInfo, err := d.Repository.FindOneChat(ctx, secret)
	if err != nil {
		return nil, err
	}

	return chatInfo, nil

}
