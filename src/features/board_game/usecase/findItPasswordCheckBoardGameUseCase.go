package usecase

import (
	"context"
	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/request"
	"time"
)

type FindItPasswordCheckBoardGameUseCase struct {
	Repository     _interface.IFindItPasswordCheckBoardGameRepository
	ContextTimeout time.Duration
}

func NewFindItPasswordCheckBoardGameUseCase(repo _interface.IFindItPasswordCheckBoardGameRepository, timeout time.Duration) _interface.IFindItPasswordCheckBoardGameUseCase {
	return &FindItPasswordCheckBoardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *FindItPasswordCheckBoardGameUseCase) FindItPasswordCheck(c context.Context, req *request.ReqFindItPasswordCheckBoardGame) (bool, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 비밀번호 확인
	passwordCheck, err := d.Repository.FindPasswordCheck(ctx, req.Password)
	if err != nil {
		return false, err
	}

	return passwordCheck, nil
}
