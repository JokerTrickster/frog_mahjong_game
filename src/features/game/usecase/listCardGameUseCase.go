package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/response"
	"time"
)

type ListCardGameUseCase struct {
	Repository     _interface.IListCardGameRepository
	ContextTimeout time.Duration
}

func NewListCardGameUseCase(repo _interface.IListCardGameRepository, timeout time.Duration) _interface.IListCardGameUseCase {
	return &ListCardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *ListCardGameUseCase) ListCard(c context.Context) (response.ResListCardGame, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	cards, err := d.Repository.FindAllBirdCard(ctx)
	if err != nil {
		return response.ResListCardGame{}, err
	}

	res := CreateResListCard(cards)
	return res, nil

}
