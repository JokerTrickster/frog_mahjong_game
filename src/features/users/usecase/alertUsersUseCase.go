package usecase

import (
	"context"
	_interface "main/features/users/model/interface"
	"main/features/users/model/request"
	"time"
)

type AlertUsersUseCase struct {
	Repository     _interface.IAlertUsersRepository
	ContextTimeout time.Duration
}

func NewAlertUsersUseCase(repo _interface.IAlertUsersRepository, timeout time.Duration) _interface.IAlertUsersUseCase {
	return &AlertUsersUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *AlertUsersUseCase) Alert(c context.Context, userID uint, req *request.ReqAlertUsers) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// userID를 찾아서 alert 활성화 여부를 변경한다.
	err := d.Repository.FindOneAndUpdateUsers(ctx, userID, req)
	if err != nil {
		return err
	}

	return nil
}
