package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"time"
)

type DiscardGameUseCase struct {
	Repository     _interface.IDiscardGameRepository
	ContextTimeout time.Duration
}

func NewDiscardGameUseCase(repo _interface.IDiscardGameRepository, timeout time.Duration) _interface.IDiscardGameUseCase {
	return &DiscardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *DiscardGameUseCase) Discard(c context.Context, userID int, req *request.ReqDiscard) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	// 플레이어가 자기 차례인지 체크
	roomUser, err := d.Repository.PlayerCheckTurn(ctx, req)
	if err != nil {
		return err
	}
	// 패를 버리기 (카드 상태 업데이트)
	err = d.Repository.UpdateCardStateDiscard(ctx, req) // 카드 상태 변경
	if err != nil {
		return err
	}
	updateRoomUser := CreateUpdateRoomUser(roomUser, req)
	// 유저 가지고 있는 패 수 -1 그리고 유저 상태 변경
	err = d.Repository.UpdateRoomUser(ctx, updateRoomUser) // 유저 상태 변경
	if err != nil {
		return err
	}

	return nil
}
