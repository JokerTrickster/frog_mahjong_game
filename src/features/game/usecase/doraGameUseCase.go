package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"time"
)

type DoraGameUseCase struct {
	Repository     _interface.IDoraGameRepository
	ContextTimeout time.Duration
}

func NewDoraGameUseCase(repo _interface.IDoraGameRepository, timeout time.Duration) _interface.IDoraGameUseCase {
	return &DoraGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *DoraGameUseCase) Dora(c context.Context, userID int, req *request.ReqDora) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	// 선플레이어가 도라를 선택했는지 체크
	err := d.Repository.CheckFirstPlayer(ctx, userID, req.RoomID)
	if err != nil {
		return err
	}
	// 카드 업데이트
	err = d.Repository.UpdateDoraCard(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
