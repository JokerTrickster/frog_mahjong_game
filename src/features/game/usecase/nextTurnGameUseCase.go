package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"time"
)

type NextTurnGameUseCase struct {
	Repository     _interface.INextTurnGameRepository
	ContextTimeout time.Duration
}

func NewNextTurnGameUseCase(repo _interface.INextTurnGameRepository, timeout time.Duration) _interface.INextTurnGameUseCase {
	return &NextTurnGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *NextTurnGameUseCase) NextTurn(c context.Context, req *request.ReqNextTurn) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	// 해당 턴 넘버를 가진 room user가 play_wait인지 확인 후 플레이 상태를 play로 변경
	err := d.Repository.UpdatePlayerNextTurn(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
