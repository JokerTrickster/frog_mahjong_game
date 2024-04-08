package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"

	"time"
)

type LoanGameUseCase struct {
	Repository     _interface.ILoanGameRepository
	ContextTimeout time.Duration
}

func NewLoanGameUseCase(repo _interface.ILoanGameRepository, timeout time.Duration) _interface.ILoanGameUseCase {
	return &LoanGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *LoanGameUseCase) Loan(c context.Context, userID uint, req *request.ReqLoan) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// loan 가능한지 체크 (마지막으로 버려진 카드인지 체크)
	err := d.Repository.CheckLoan(ctx, req)
	if err != nil {
		return err
	}

	// loan 하기 (상대방이 버린 카드를 가져온다)
	err = d.Repository.Loan(ctx, req)
	if err != nil {
		return err
	}

	// 룸 유저 카드 수와 상태값을 변경한다.
	err = d.Repository.UpdateRoomUserCardCount(ctx, userID, uint(req.RoomID))
	if err != nil {
		return err
	}
	return nil
}
