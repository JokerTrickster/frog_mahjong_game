package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"time"
)

type WinRequestGameUseCase struct {
	Repository     _interface.IWinRequestGameRepository
	ContextTimeout time.Duration
}

func NewWinRequestGameUseCase(repo _interface.IWinRequestGameRepository, timeout time.Duration) _interface.IWinRequestGameUseCase {
	return &WinRequestGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *WinRequestGameUseCase) WinRequest(c context.Context, req *request.ReqWinRequest) (bool, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// userID와 roomID를 통해 room User 정보를 가져온다.
	roomUser, err := d.Repository.GetRoomUser(ctx, req.UserID, req.RoomID)
	if err != nil {
		return false, err
	}

	// 현재 플레이 상태가 play or loan 인지 체크 및 카드 수가 6장인지 체크 후 5점이상이면 true, 아니면 false
	if IsCheckedWinRequest(roomUser) {
		return true, nil
	}

	return false, nil
}
