package usecase

import (
	"context"
	_interface "main/features/users/model/interface"
	"time"
)

type OneCoinUsersUseCase struct {
	Repository     _interface.IOneCoinUsersRepository
	ContextTimeout time.Duration
}

func NewOneCoinUsersUseCase(repo _interface.IOneCoinUsersRepository, timeout time.Duration) _interface.IOneCoinUsersUseCase {
	return &OneCoinUsersUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *OneCoinUsersUseCase) OneCoin(c context.Context) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 유저 정보를 변경한다.
	err := d.Repository.OneCoin(ctx)
	if err != nil {
		return err
	}

	return nil
}
