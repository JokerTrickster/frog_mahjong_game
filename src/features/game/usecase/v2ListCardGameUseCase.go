package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/response"
	"time"
)

type V2ListCardGameUseCase struct {
	Repository     _interface.IV2ListCardGameRepository
	ContextTimeout time.Duration
}

func NewV2ListCardGameUseCase(repo _interface.IV2ListCardGameRepository, timeout time.Duration) _interface.IV2ListCardGameUseCase {
	return &V2ListCardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *V2ListCardGameUseCase) V2ListCard(c context.Context) (response.ResV2ListCardGame, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	cards, err := d.Repository.FindAllBirdCard(ctx)
	if err != nil {
		return response.ResV2ListCardGame{}, err
	}
	count, err := d.Repository.CountAllBirdCard(ctx)
	if err != nil {
		return response.ResV2ListCardGame{}, err
	}

	res := CreateResV2ListCard(cards,count)
	return res, nil

}
