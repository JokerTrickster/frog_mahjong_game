package usecase

import (
	"context"
	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/request"
	"time"
)

type FindItCoinBoardGameUseCase struct {
	Repository     _interface.IFindItCoinBoardGameRepository
	ContextTimeout time.Duration
}

func NewFindItCoinBoardGameUseCase(repo _interface.IFindItCoinBoardGameRepository, timeout time.Duration) _interface.IFindItCoinBoardGameUseCase {
	return &FindItCoinBoardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *FindItCoinBoardGameUseCase) FindItCoin(c context.Context, userID int, req *request.ReqFindItCoinBoardGame) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	err := d.Repository.UpdateUserCoin(ctx, userID, req.Coin)
	if err != nil {
		return err
	}

	return nil
}
