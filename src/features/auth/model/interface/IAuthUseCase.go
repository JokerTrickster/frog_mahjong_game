package _interface

import (
	"context"
	"main/features/auth/model/request"
	"main/features/auth/model/response"
)

type ISignupAuthUseCase interface {
	Signup(c context.Context, req *request.ReqSignup) error
}

type ISigninAuthUseCase interface {
	Signin(c context.Context, req *request.ReqSignin) (response.ResSignin, error)
}

type ILogoutAuthUseCase interface {
	Logout(c context.Context, uID uint) error
}

type IReissueAuthUseCase interface {
	Reissue(c context.Context, req *request.ReqReissue) error
}
