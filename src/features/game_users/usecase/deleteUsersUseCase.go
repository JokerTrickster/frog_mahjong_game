package usecase

import (
	"context"
	_interface "main/features/game_users/model/interface"
	"time"
)

type DeleteUsersUseCase struct {
	Repository     _interface.IDeleteUsersRepository
	ContextTimeout time.Duration
}

func NewDeleteUsersUseCase(repo _interface.IDeleteUsersRepository, timeout time.Duration) _interface.IDeleteUsersUseCase {
	return &DeleteUsersUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *DeleteUsersUseCase) Delete(c context.Context, userID uint) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 유저 정보를 삭제한다.
	err := d.Repository.FindOneAndDeleteUsers(ctx, userID)
	if err != nil {
		return err
	}

	// 토큰 제거
	err = d.Repository.DeleteToken(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
