package usecase

import (
	"context"
	"fmt"
	_interface "main/features/auth/model/interface"
	"main/features/auth/model/request"
	"time"
)

type ReissueAuthUseCase struct {
	Repository     _interface.IReissueAuthRepository
	ContextTimeout time.Duration
}

func NewReissueAuthUseCase(repo _interface.IReissueAuthRepository, timeout time.Duration) _interface.IReissueAuthUseCase {
	return &ReissueAuthUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *ReissueAuthUseCase) Reissue(c context.Context, req *request.ReqReissue) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	fmt.Println(ctx)
	return nil
}
