package usecase

import (
	"context"
	_interface "main/features/auth/model/interface"
	"main/features/auth/model/request"
	"main/features/auth/model/response"
	"time"
)

type CheckSigninAuthUseCase struct {
	Repository     _interface.ICheckSigninAuthRepository
	ContextTimeout time.Duration
}

func NewCheckSigninAuthUseCase(repo _interface.ICheckSigninAuthRepository, timeout time.Duration) _interface.ICheckSigninAuthUseCase {
	return &CheckSigninAuthUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *CheckSigninAuthUseCase) CheckSignin(c context.Context, req *request.ReqCheckSignin) (bool, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	// user check
	user, err := d.Repository.FindOneAndUpdateUser(ctx, req.Email, req.Password)
	if err != nil {
		return false, err
	}

	// 기존 토큰이 있는지 체크
	prevTokens, err := d.Repository.CheckToken(ctx, user.ID)
	if err != nil {
		return false, err
	}
	res := response.ResSignin{
		IsDuplicateLogin: false,
	}
	if prevTokens != nil {
		res.IsDuplicateLogin = true
	}
	if prevTokens != nil {
		return true, nil
	}
	return false, nil
}
