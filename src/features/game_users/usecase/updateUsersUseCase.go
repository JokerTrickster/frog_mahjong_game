package usecase

import (
	"context"
	_interface "main/features/game_users/model/interface"
	"main/features/game_users/model/request"
	"time"
)

type UpdateUsersUseCase struct {
	Repository     _interface.IUpdateUsersRepository
	ContextTimeout time.Duration
}

func NewUpdateUsersUseCase(repo _interface.IUpdateUsersRepository, timeout time.Duration) _interface.IUpdateUsersUseCase {
	return &UpdateUsersUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *UpdateUsersUseCase) Update(c context.Context, userID uint, req *request.ReqUpdateGameUsers) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// entitySQL 을 생성한다.
	entitySQL := CreateUpdateUsersEntitySQL(userID, req)

	// 유저 정보를 변경한다.
	err := d.Repository.FindOneAndUpdateUsers(ctx, &entitySQL)
	if err != nil {
		return err
	}

	return nil
}
