package usecase

import (
	"context"
	"fmt"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"time"
)

type NextTurnGameUseCase struct {
	Repository     _interface.INextTurnGameRepository
	ContextTimeout time.Duration
}

func NewNextTurnGameUseCase(repo _interface.INextTurnGameRepository, timeout time.Duration) _interface.INextTurnGameUseCase {
	return &NextTurnGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *NextTurnGameUseCase) NextTurn(c context.Context, req *request.ReqNextTurn) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	fmt.Println("NextTurnGameUseCase NextTurn", ctx)

	return nil
}
