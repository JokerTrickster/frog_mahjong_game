package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"time"
)

type UpdateCardGameUseCase struct {
	Repository     _interface.IUpdateCardGameRepository
	ContextTimeout time.Duration
}

func NewUpdateCardGameUseCase(repo _interface.IUpdateCardGameRepository, timeout time.Duration) _interface.IUpdateCardGameUseCase {
	return &UpdateCardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *UpdateCardGameUseCase) UpdateCard(c context.Context, req *request.ReqUpdateCard) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// bird cards DTO 생성
	birdCardsDTO := UpdateBirdCardsDTO(req)

	// db 저장
	err := d.Repository.UpdateCard(ctx, birdCardsDTO)
	if err != nil {
		return err
	}

	return nil

}
