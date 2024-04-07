package usecase

import (
	"context"
	"fmt"
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

func (d *LoanGameUseCase) Loan(c context.Context, req *request.ReqLoan) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	fmt.Println(ctx)

	return nil
}
