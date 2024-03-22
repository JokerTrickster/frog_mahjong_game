package usecase

import (
	"context"
	_interface "main/features/auth/model/interface"
	"main/features/auth/model/request"
	"main/features/auth/model/response"
	"main/utils"
	"time"
)

type SigninAuthUseCase struct {
	Repository     _interface.ISigninAuthRepository
	ContextTimeout time.Duration
}

func NewSigninAuthUseCase(repo _interface.ISigninAuthRepository, timeout time.Duration) _interface.ISigninAuthUseCase {
	return &SigninAuthUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *SigninAuthUseCase) Signin(c context.Context, req *request.ReqSignin) (response.ResSignin, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	// user check
	user, err := d.Repository.FindOneAndUpdateUser(ctx, req.Email, req.Password)
	if err != nil {
		return response.ResSignin{}, err
	}

	// token create
	accessToken, refreshToken, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		return response.ResSignin{}, err
	}

	//response create
	res := response.ResSignin{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}
