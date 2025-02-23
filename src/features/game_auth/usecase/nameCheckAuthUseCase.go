package usecase

import (
	"context"
	_interface "main/features/game_auth/model/interface"
	"main/features/game_auth/model/request"
	"time"
)

type NameCheckAuthUseCase struct {
	Repository     _interface.INameCheckAuthRepository
	ContextTimeout time.Duration
}

func NewNameCheckAuthUseCase(repo _interface.INameCheckAuthRepository, timeout time.Duration) _interface.INameCheckAuthUseCase {
	return &NameCheckAuthUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *NameCheckAuthUseCase) NameCheck(c context.Context, req *request.ReqGameNameCheck) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 이름 중복 체크
	err := d.Repository.CheckName(ctx, req.Name)
	if err != nil {
		return err
	}
	return nil
}
