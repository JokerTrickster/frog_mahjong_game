package usecase

import (
	"context"
	_interface "main/features/game_auth/model/interface"
	"main/features/game_auth/model/request"
	"time"
)

type FCMTokenAuthUseCase struct {
	Repository     _interface.IFCMTokenAuthRepository
	ContextTimeout time.Duration
}

func NewFCMTokenAuthUseCase(repo _interface.IFCMTokenAuthRepository, timeout time.Duration) _interface.IFCMTokenAuthUseCase {
	return &FCMTokenAuthUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *FCMTokenAuthUseCase) FCMToken(c context.Context, userID uint, req *request.ReqGameFCMToken) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// userID를 찾아서 FCM Token을 저장한다.
	err := d.Repository.SaveFCMToken(ctx, userID, req.FCMToken)
	if err != nil {
		return err
	}

	return nil
}
