package usecase

import (
	"context"
	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/response"
	"time"
)

type FrogCardListBoardGameUseCase struct {
	Repository     _interface.IFrogCardListBoardGameRepository
	ContextTimeout time.Duration
}

func NewFrogCardListBoardGameUseCase(repo _interface.IFrogCardListBoardGameRepository, timeout time.Duration) _interface.IFrogCardListBoardGameUseCase {
	return &FrogCardListBoardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *FrogCardListBoardGameUseCase) FrogCardList(c context.Context) (response.ResFrogCardListBoardGame, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	cards, err := d.Repository.FindFrogCard(ctx)
	if err != nil {
		return response.ResFrogCardListBoardGame{}, err
	}
	count, err := d.Repository.CountFrogCard(ctx)
	if err != nil {
		return response.ResFrogCardListBoardGame{}, err
	}

	res := CreateResFrogCardList(cards, count)
	return res, nil

}
