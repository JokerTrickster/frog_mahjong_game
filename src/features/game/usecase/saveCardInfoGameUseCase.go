package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"time"
)

type SaveCardInfoGameUseCase struct {
	Repository     _interface.ISaveCardInfoGameRepository
	ContextTimeout time.Duration
}

func NewSaveCardInfoGameUseCase(repo _interface.ISaveCardInfoGameRepository, timeout time.Duration) _interface.ISaveCardInfoGameUseCase {
	return &SaveCardInfoGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *SaveCardInfoGameUseCase) SaveCardInfo(c context.Context, req *request.ReqSaveCardInfo) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// bird cards DTO 생성
	birdCardsDTO := CreateBirdCardsDTO(req)
	
	// db 저장
	err := d.Repository.SaveCardInfo(ctx, birdCardsDTO)
	if err != nil {
		return err
	}

	return nil

}
