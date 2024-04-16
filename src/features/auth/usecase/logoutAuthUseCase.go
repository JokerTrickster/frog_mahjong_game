package usecase

import (
	"context"
	_interface "main/features/auth/model/interface"
	"time"
)

type LogoutAuthUseCase struct {
	Repository        _interface.ILogoutAuthRepository
	ContextTimeLogout time.Duration
}

func NewLogoutAuthUseCase(repo _interface.ILogoutAuthRepository, timeLogout time.Duration) _interface.ILogoutAuthUseCase {
	return &LogoutAuthUseCase{Repository: repo, ContextTimeLogout: timeLogout}
}

func (d *LogoutAuthUseCase) Logout(c context.Context, uID uint) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeLogout)
	defer cancel()
	err := d.Repository.FindOneAndUpdateUser(ctx, uID)
	if err != nil {
		return err
	}

	return nil
}
