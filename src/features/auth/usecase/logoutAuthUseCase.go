package usecase

import (
	"context"
	_interface "main/features/auth/model/interface"
	"time"
)

type LogoutAuthUseCase struct {
	Repository     _interface.ILogoutAuthRepository
	ContextTimeout time.Duration
}

func NewLogoutAuthUseCase(repo _interface.ILogoutAuthRepository, timeout time.Duration) _interface.ILogoutAuthUseCase {
	return &LogoutAuthUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *LogoutAuthUseCase) Logout(c context.Context, uID uint) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	err := d.Repository.FindOneAndUpdateUser(ctx, uID)
	if err != nil {
		return err
	}

	// 토큰 제거
	err = d.Repository.DeleteToken(ctx, uID)
	if err != nil {
		return err
	}

	return nil
}
