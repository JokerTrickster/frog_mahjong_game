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
