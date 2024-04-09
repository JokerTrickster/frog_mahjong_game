package usecase

import (
	"context"
	"fmt"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"time"
)

type WinRequestGameUseCase struct {
	Repository     _interface.IWinRequestGameRepository
	ContextTimeout time.Duration
}

func NewWinRequestGameUseCase(repo _interface.IWinRequestGameRepository, timeout time.Duration) _interface.IWinRequestGameUseCase {
	return &WinRequestGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *WinRequestGameUseCase) WinRequest(c context.Context, req *request.ReqWinRequest) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	fmt.Println(ctx)

	return nil
}
