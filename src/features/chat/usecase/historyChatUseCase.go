package usecase

import (
	"context"
	_interface "main/features/chat/model/interface"
	"main/features/chat/model/request"
	"main/features/chat/model/response"
	"time"
)

type HistoryChatUseCase struct {
	Repository     _interface.IHistoryChatRepository
	ContextTimeout time.Duration
}

func NewHistoryChatUseCase(repo _interface.IHistoryChatRepository, timeout time.Duration) _interface.IHistoryChatUseCase {
	return &HistoryChatUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *HistoryChatUseCase) History(c context.Context, req *request.ReqHistory) (response.ResHistoryChat, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	entitySQL := CreateChatHistoryEntitySQL(req)
	chats, err := d.Repository.FindChatHistory(ctx, entitySQL)
	if err != nil {
		return response.ResHistoryChat{}, err
	}
	total, err := d.Repository.CountChatHistory(ctx, entitySQL)
	if err != nil {
		return response.ResHistoryChat{}, err
	}

	//create res
	res, err := CreateResHistoryChat(chats, total)
	if err != nil {
		return response.ResHistoryChat{}, err
	}

	return res, nil

}
