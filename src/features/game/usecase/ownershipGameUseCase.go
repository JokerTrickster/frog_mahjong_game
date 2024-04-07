package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"time"
)

type OwnershipGameUseCase struct {
	Repository     _interface.IOwnershipGameRepository
	ContextTimeout time.Duration
}

func NewOwnershipGameUseCase(repo _interface.IOwnershipGameRepository, timeout time.Duration) _interface.IOwnershipGameUseCase {
	return &OwnershipGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *OwnershipGameUseCase) Ownership(c context.Context, req *request.ReqOwnership) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	// 카드 상태 없데이트
	err := d.Repository.UpdateCardState(ctx, req)
	if err != nil {
		return err
	}
	// 소유 카드 수 업데이트
	// 유저id로 room_users에서 찾아서 card_count를 더한 후 업데이트 한다.
	err = d.Repository.UpdateRoomUserCardCount(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
