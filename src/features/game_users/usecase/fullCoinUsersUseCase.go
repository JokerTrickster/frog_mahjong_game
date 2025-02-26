package usecase

import (
	"context"
	_interface "main/features/game_users/model/interface"
	"time"
)

type FullCoinUsersUseCase struct {
	Repository     _interface.IFullCoinUsersRepository
	ContextTimeout time.Duration
}

func NewFullCoinUsersUseCase(repo _interface.IFullCoinUsersRepository, timeout time.Duration) _interface.IFullCoinUsersUseCase {
	return &FullCoinUsersUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *FullCoinUsersUseCase) FullCoin(c context.Context) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 유저 정보를 변경한다.
	err := d.Repository.FullCoin(ctx)
	if err != nil {
		return err
	}

	return nil
}
